package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"export-service/internal/core/domain"
	"export-service/internal/core/ports"
	"export-service/internal/repositories/crm_company_repo"
	"export-service/internal/repositories/crm_solicitation_repo"
	"export-service/internal/services/crm_exporter"
	"export-service/internal/services/data_presenter"

	"go.uber.org/zap"
)

type CrmExportUseCase struct {
	downloader           ports.Downloader
	presentationSpecRepo ports.PresentationSpecRepository
	companyRepo          *crm_company_repo.PgCrmCompanyRepository
	solicitationRepo     *crm_solicitation_repo.PgCrmSolicitationRepository
	mailer               ports.Mailer
	logger               *zap.Logger
}

func NewCrmExportUseCase(downloader ports.Downloader, presentationSpecRepo ports.PresentationSpecRepository, companyRepo *crm_company_repo.PgCrmCompanyRepository, solicitationRepo *crm_solicitation_repo.PgCrmSolicitationRepository, mailer ports.Mailer, logger *zap.Logger) *CrmExportUseCase {
	return &CrmExportUseCase{
		downloader:           downloader,
		presentationSpecRepo: presentationSpecRepo,
		companyRepo:          companyRepo,
		solicitationRepo:     solicitationRepo,
		mailer:               mailer,
		logger:               logger,
	}
}

func (c *CrmExportUseCase) Execute(request CrmExportRequest, requestConfigs map[string]any) error {

	crm, ok := requestConfigs["crm"].(string)
	if !ok {
		return errors.New("invalid or missing crm header")
	}

	solicitation, err := c.solicitationRepo.Create(context.Background(), crm_solicitation_repo.CreateSolicitation{
		ListId:        request.ListID,
		UserEmail:     request.UserEmail,
		Current:       0,
		Total:         int(requestConfigs["total"].(int64)),
		OwnerId:       requestConfigs["owner_id"].(string),
		PipelineId:    requestConfigs["pipeline_id"].(string),
		StageId:       requestConfigs["stage_id"].(string),
		OverwriteData: requestConfigs["overwrite_data"].(bool),
		CreateDeal:    requestConfigs["create_deal"].(bool),
	})
	if err != nil {
		c.logError("Error creating solicitation in db", err, request)
		return err
	}

	data, err := c.downloadData(request)
	if err != nil {
		c.logError("Error when downloading data", err, request)
		return err
	}

	spec, err := c.getPresentationSpec(request, crm)
	if err != nil {
		c.logError("Error when getting presentation spec", err, request)
		return err
	}

	presentedData, err := c.applyPresentationSpecCrm(request, data, spec)
	if err != nil {
		c.logError("Error when applying presentation spec", err, request)
		return err
	}

	crmService, exists := crm_exporter.GetCrm(crm, c.companyRepo)
	if !exists {
		c.logError("Crm is invalid", err, request)
		return err
	}

	crmClient, err := crmService.Authorize(context.Background(), request.UserCompany)
	if err != nil {
		c.logError("Error when authorizing crm client", err, request)
		return err
	}

	presentedDataMappedToCnpjs, err := c.mapPresentedDataToCnpjs(data, presentedData)
	if err != nil {
		c.logError("Error mapping cnpj to presented data", err, request)
		return err
	}

	err = c.sendAllLeads(crmService, crmClient, presentedDataMappedToCnpjs, requestConfigs, solicitation.ListId)

	return err

	// err = c.sendEmail(request, url)
	// if err != nil {
	// 	c.logError("Error when sending email", err, request)
	// 	return "", err
	// }

	// return url, nil
}

func (c *CrmExportUseCase) mapPresentedDataToCnpjs(data []map[string]any, presentedData []map[string]any) (map[any]map[string]any, error) {
	presentedDataCnpjMap := make(map[any]map[string]any)
	for key, item := range presentedData {
		if len(data) != len(presentedData) {
			return map[any]map[string]any{}, errors.New("data and presented data length mismatch")
		}
		cnpj, exists := data[key]["cnpj"]
		if !exists {
			return map[any]map[string]any{}, errors.New("data and presented data length mismatch")
		}

		presentedDataCnpjMap[cnpj] = item
	}

	return presentedDataCnpjMap, nil
}

func (c *CrmExportUseCase) sendAllLeads(crmService crm_exporter.Crm, client any, leadsData map[any]map[string]any, configs map[string]any, solicitationId string) error {
	for cnpj, leadData := range leadsData {
		leadResult, err := crmService.SendLead(client, leadData, configs)
		c.updateExportedCompaniesInDb(leadResult, cnpj, solicitationId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CrmExportUseCase) updateExportedCompaniesInDb(leadResult crm_exporter.CreatedLead, cnpj any, solicitationId string) error {

	c.solicitationRepo.Update(context.Background(), crm_solicitation_repo.UpdateExportedCompaniesParms{
		Cnpj:               cnpj,
		NewExportedCompany: leadResult,
	}, solicitationId)

	return nil
}

func (c *CrmExportUseCase) downloadData(request CrmExportRequest) ([]map[string]any, error) {
	c.logInfo("Downloading data", request)
	bytes, err := c.downloader.Download(request.DataDownloadURL)
	if err != nil {
		return nil, err
	}

	var data []map[string]any
	err = json.Unmarshal(bytes, &data)
	return data, err
}

func (c *CrmExportUseCase) getPresentationSpec(request CrmExportRequest, crm string) (domain.PresentationSpec, error) {
	c.logInfo("Getting presentation spec", request)

	return c.presentationSpecRepo.Get(context.Background(), ports.PresentationSpecQueryParams{
		UserEmail:   request.UserEmail,
		UserCompany: request.UserCompany,
		Service:     "crm_" + crm,
		DataSource:  request.DataSource,
	})
}

func (c *CrmExportUseCase) applyPresentationSpecCrm(request CrmExportRequest, data []map[string]any, spec domain.PresentationSpec) ([]map[string]any, error) {
	c.logInfo("Applying presentation spec", request)
	result := make([]map[string]any, 0)

	for _, d := range data {
		dd, err := data_presenter.PresentSingle(d, spec.Spec)
		if err != nil {
			return nil, err
		}
		result = append(result, dd)
	}

	return result, nil
}

func (c *CrmExportUseCase) sendEmail(request CrmExportRequest, url string) error {
	c.logInfo("Sending email", request)
	return c.mailer.SendEmail(request.UserEmail, request.UserName, "d-b1c0014b01eb410a8c4b8112e4418a3f", url)
}

func (c *CrmExportUseCase) logInfo(message string, request CrmExportRequest) {
	c.logger.Info(message, zap.Any("request", request))
}

func (c *CrmExportUseCase) logError(message string, err error, request CrmExportRequest) {
	c.logger.Error(message, zap.Error(err), zap.Any("request", request))
}
