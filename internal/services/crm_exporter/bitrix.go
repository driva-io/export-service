package crm_exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"export-service/internal/core/ports"
	"export-service/internal/repositories/crm_company_repo"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type BitrixClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewBitrixClient(baseURL string) *BitrixClient {
	return &BitrixClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (bc *BitrixClient) MakeRequest(method, endpoint string, body any) (map[string]any, error) {
	url := fmt.Sprintf("%s%s", bc.BaseURL, endpoint)

	var jsonBody []byte
	var err error
	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bitrix http error: %s", responseBody)
	}

	// Parse the response body into a map
	var result map[string]any
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return result, nil
}

type BitrixService struct {
	companyRepo *crm_company_repo.PgCrmCompanyRepository
}

func NewBitrixService(companyRepo *crm_company_repo.PgCrmCompanyRepository) *BitrixService {
	return &BitrixService{companyRepo: companyRepo}
}

func (b *BitrixService) Authorize(ctx context.Context, workspaceId string) (any, error) {
	company, err := b.companyRepo.GetCompanyByWorkspaceId(ctx, ports.CrmCompanyQueryParams{Crm: "bitrix", WorkspaceId: workspaceId})
	if err != nil {
		return nil, err
	}

	return NewBitrixClient(company.Token.String), nil
}

func (b *BitrixService) Validate(ctx *fiber.Ctx, client any) bool {
	bitrixClient, ok := client.(*BitrixClient)
	if !ok {
		return false
	}

	_, err := bitrixClient.MakeRequest("GET", "user.get", nil)
	return err == nil
}

func (b *BitrixService) Install(installData any) (any, error) {
	// TODO: Implement the method
	return nil, nil
}

func (b *BitrixService) OAuthCallback(ctx *fiber.Ctx, params ...any) (any, error) {
	panic("unimplemented")
}

func (b *BitrixService) SendLead(client any, mappedStorageData map[string]any, correspondingRawData map[string]any, configs map[string]any, existingLead map[string]any) (CreatedLead, error) {
	bitrixClient, ok := client.(*BitrixClient)
	if !ok {
		return CreatedLead{}, errors.New("invalid Bitrix client")
	}

	workspaceId, err := getConfigValue[string](configs, "workspace_id")
	if err != nil {
		return CreatedLead{}, err
	}

	companyConfigs, err := b.companyRepo.GetCompanyByWorkspaceId(context.Background(), ports.CrmCompanyQueryParams{Crm: "bitrix", WorkspaceId: workspaceId})
	if err != nil {
		return CreatedLead{}, err
	}

	ownerId, err := getConfigValue[string](configs, "owner_id")
	if err != nil {
		return CreatedLead{}, err
	}

	stageId, err := getConfigValue[string](configs, "stage_id")
	if err != nil {
		return CreatedLead{}, err
	}

	pipelineId, err := getConfigValue[string](configs, "pipeline_id")
	if err != nil {
		return CreatedLead{}, err
	}

	leadFormat := false
	if companyConfigs.Environment.String == "classic" {
		leadFormat = true
	}
	createDeal := configs["create_deal"].(bool)
	createdLead := CreatedLead{}

	if company, exists := mappedStorageData["company"]; exists {
		companyStatus, err := processBitrixCompany(bitrixClient, company, existingLead, correspondingRawData, ownerId)
		if err != nil {
			return createdLead, err
		}
		createdLead.Company = companyStatus
	}

	if deal, exists := mappedStorageData["deal"]; exists && createDeal && !leadFormat {
		dealStatus, err := processBitrixDeal(bitrixClient, deal, existingLead, correspondingRawData, ownerId, pipelineId, stageId)
		if err != nil {
			return createdLead, err
		}
		createdLead.Deal = dealStatus
	}

	if contact, exists := mappedStorageData["contact"]; exists {
		contactStatus, err := processBitrixContact(bitrixClient, contact, existingLead, correspondingRawData, ownerId)
		if err != nil {
			return createdLead, err
		}
		createdLead.Contacts = contactStatus
	}

	if contacts, exists := mappedStorageData["contacts"]; exists {
		contactsStatus, err := processBitrixContacts(bitrixClient, contacts, existingLead, correspondingRawData, ownerId)
		if err != nil {
			return createdLead, err
		}
		createdLead.Contacts = contactsStatus
	}

	if lead, exists := mappedStorageData["lead"]; exists && leadFormat {
		leadStatus, err := processBitrixLead(bitrixClient, lead, existingLead, correspondingRawData, ownerId)
		if err != nil {
			return createdLead, err
		}
		createdLead.Lead = leadStatus
	}

	createdLead, err = createBitrixAssociations(bitrixClient, createdLead)
	return createdLead, err
}

func createBitrixAssociations(client *BitrixClient, lead CreatedLead) (CreatedLead, error) {

	if lead.Deal != nil && lead.Company != nil {
		if !associationExists(lead.Deal.Associations, "company", lead.Company.CrmId) {
			_, err := client.MakeRequest("POST", "crm.deal.update", map[string]any{
				"id":     lead.Deal.CrmId,
				"fields": map[string]any{"COMPANY_ID": lead.Company.CrmId},
			})
			if err != nil {
				return lead, err
			}

			lead.Deal.Associations = append(lead.Deal.Associations, Association{
				ObjectType: "company",
				CrmId:      lead.Company.CrmId,
			})
		}
	}

	if lead.Deal != nil && lead.Contacts != nil {
		var contactIds []any
		for _, contact := range *lead.Contacts {
			contactIds = append(contactIds, contact.CrmId)
		}

		if !associationExists(lead.Deal.Associations, "contacts", contactIds) {
			_, err := client.MakeRequest("POST", "crm.deal.update", map[string]any{
				"id":     lead.Deal.CrmId,
				"fields": map[string]any{"CONTACT_IDS": contactIds},
			})
			if err != nil {
				return lead, err
			}

			lead.Deal.Associations = append(lead.Deal.Associations, Association{
				ObjectType: "contacts",
				CrmId:      contactIds,
			})
		}
	}

	if lead.Company != nil && lead.Contacts != nil {
		for i := range *lead.Contacts {
			contact := &(*lead.Contacts)[i]

			if !associationExists(contact.Associations, "company", lead.Company.CrmId) {
				_, err := client.MakeRequest("POST", "crm.contact.update", map[string]any{
					"id":     contact.CrmId,
					"fields": map[string]any{"COMPANY_ID": lead.Company.CrmId},
				})
				if err != nil {
					return lead, err
				}

				contact.Associations = append(contact.Associations, Association{
					ObjectType: "company",
					CrmId:      lead.Company.CrmId,
				})
			}
		}
	}

	if lead.Company != nil && lead.Lead != nil {
		if !associationExists(lead.Lead.Associations, "company", lead.Company.CrmId) {
			_, err := client.MakeRequest("POST", "crm.lead.update", map[string]any{
				"id":     lead.Lead.CrmId,
				"fields": map[string]any{"COMPANY_ID": lead.Company.CrmId},
			})
			if err != nil {
				return lead, err
			}

			lead.Lead.Associations = append(lead.Lead.Associations, Association{
				ObjectType: "company",
				CrmId:      lead.Company.CrmId,
			})
		}
	}

	if lead.Lead != nil && lead.Contacts != nil {
		var contactIds []any
		for _, contact := range *lead.Contacts {
			contactIds = append(contactIds, contact.CrmId)
		}

		if !associationExists(lead.Lead.Associations, "contacts", contactIds) {
			_, err := client.MakeRequest("POST", "crm.lead.update", map[string]any{
				"id":     lead.Lead.CrmId,
				"fields": map[string]any{"CONTACT_IDS": contactIds},
			})
			if err != nil {
				return lead, err
			}

			lead.Lead.Associations = append(lead.Lead.Associations, Association{
				ObjectType: "contacts",
				CrmId:      contactIds,
			})
		}
	}

	return lead, nil
}

func associationExists(associations []Association, objectType string, crmId any) bool {
	switch v := crmId.(type) {
	case []any:
		for _, id := range v {
			for _, assoc := range associations {
				if assoc.ObjectType == objectType {
					switch assocCrmId := assoc.CrmId.(type) {
					case []any:
						if reflect.DeepEqual(v, assocCrmId) {
							return true
						}
					default:
						if assoc.CrmId == id {
							return true
						}
					}
				}
			}
		}
	default:
		for _, assoc := range associations {
			if assoc.ObjectType == objectType {
				switch assocCrmId := assoc.CrmId.(type) {
				case []any:
					if reflect.DeepEqual([]any{crmId}, assocCrmId) {
						return true
					}
				default:
					if assoc.CrmId == crmId {
						return true
					}
				}
			}
		}
	}

	return false
}

func (b *BitrixService) GetPipelines(client any) ([]Pipeline, error) {
	// TODO: Implement the method
	return nil, nil
}

func (b *BitrixService) GetFields(client any) (CrmFields, error) {
	// TODO: Implement the method
	return CrmFields{}, nil
}

func (b *BitrixService) GetOwners(client any) ([]Owner, error) {
	// TODO: Implement the method
	return nil, nil
}

func processBitrixCompany(client *BitrixClient, company any, existingLead, rawData map[string]any, ownerId string) (*ObjectStatus, error) {
	exportedCompany, exists := existingLead["company"].(map[string]any)
	if exists && exportedCompany["crm_id"] != nil {
		return createExistingStatus(exportedCompany), nil
	}

	companyData, ok := company.(map[string]any)
	if !ok {
		return nil, errors.New("invalid company data")
	}

	sentCompany, err := sendBitrixCompany(client, companyData, ownerId)
	if err != nil {
		return nil, err
	}

	if drivaID, exists := rawData["company_contact_id"].(string); exists {
		sentCompany.DrivaContactId = drivaID
	}

	return &sentCompany, nil
}

func sendBitrixCompany(client *BitrixClient, mappedCompanyData map[string]any, ownerId string) (ObjectStatus, error) {
	companyEntity, exists := mappedCompanyData["entity"]
	if !exists {
		return ObjectStatus{
			Status:  Failed,
			Message: "company entity not found in mapped company data",
		}, errors.New("company entity not found in mapped company data")
	}
	companyEntityMap, isMap := companyEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{
			Status:  Failed,
			Message: "company entity is not a map",
		}, errors.New("company entity is not a map")
	}

	companyEntityMap["ASSIGNED_BY_ID"] = ownerId

	searchFilters := map[string]any{"TITLE": companyEntityMap["TITLE"]}
	existingCompany, err := searchForExistingBitrixObject(client, "company", searchFilters)
	if err != nil {
		return ObjectStatus{
			Status:  Failed,
			Message: err.Error(),
		}, err
	}

	searchFields := buildSearchFields(searchFilters)

	var company map[string]any
	var status Status
	var message string
	if existingCompany != nil {
		updatedCompany, err := client.MakeRequest("POST", "crm.company.update", map[string]any{"id": existingCompany.(map[string]any)["ID"], "fields": companyEntityMap})
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		updatedCompany["result"] = existingCompany.(map[string]any)["ID"]
		status = Updated
		company = updatedCompany
		message = fmt.Sprintf("Searched fields: %s", searchFields)
	} else {
		createdCompany, err := client.MakeRequest("POST", "crm.company.add", map[string]any{"fields": companyEntityMap})
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Created
		company = createdCompany
	}

	return ObjectStatus{
		CrmId:   company["result"],
		Status:  status,
		Message: message,
	}, nil
}

func searchForExistingBitrixObject(client *BitrixClient, objectType string, searchFilters map[string]any) (any, error) {
	existingObject, err := client.MakeRequest("POST", "crm."+objectType+".list", map[string]any{"filter": searchFilters})

	if err != nil {
		return nil, err
	}

	if len(existingObject["result"].([]any)) > 0 {
		return existingObject["result"].([]any)[0], nil
	}

	return nil, nil
}

func processBitrixDeal(client *BitrixClient, deal any, existingLead, rawData map[string]any, ownerId, pipelineId, stageId string) (*ObjectStatus, error) {
	exportedDeal, exists := existingLead["deal"].(map[string]any)
	if exists && exportedDeal["crm_id"] != nil {
		return createExistingStatus(exportedDeal), nil
	}

	dealData, ok := deal.(map[string]any)
	if !ok {
		return nil, errors.New("invalid deal data")
	}

	sentDeal, err := sendBitrixDeal(client, dealData, ownerId, pipelineId, stageId)
	if err != nil {
		return nil, err
	}

	if drivaID, exists := rawData["company_contact_id"].(string); exists {
		sentDeal.DrivaContactId = drivaID
	}

	return &sentDeal, nil
}

func sendBitrixDeal(client *BitrixClient, mappedDealData map[string]any, ownerId, pipelineId, stageId string) (ObjectStatus, error) {
	dealEntity, exists := mappedDealData["entity"]
	if !exists {
		return ObjectStatus{
			Status:  Failed,
			Message: "deal entity not found in mapped deal data",
		}, errors.New("deal entity not found in mapped deal data")
	}
	dealEntityMap, isMap := dealEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{
			Status:  Failed,
			Message: "deal entity is not a map",
		}, errors.New("deal entity is not a map")
	}

	dealEntityMap["CATEGORY_ID"] = pipelineId
	dealEntityMap["STAGE_ID"] = stageId
	dealEntityMap["ASSIGNED_BY_ID"] = ownerId

	searchFilters := map[string]any{"TITLE": dealEntityMap["TITLE"]}
	existingDeal, err := searchForExistingBitrixObject(client, "deal", searchFilters)
	if err != nil {
		return ObjectStatus{
			Status:  Failed,
			Message: err.Error(),
		}, err
	}

	searchFields := buildSearchFields(searchFilters)

	var deal map[string]any
	var status Status
	var message string
	if existingDeal != nil {
		updatedDeal, err := client.MakeRequest("POST", "crm.deal.update", map[string]any{"id": existingDeal.(map[string]any)["ID"], "fields": dealEntityMap})
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		updatedDeal["result"] = existingDeal.(map[string]any)["ID"]
		status = Updated
		deal = updatedDeal
		message = fmt.Sprintf("Searched fields: %s", searchFields)
	} else {
		createdDeal, err := client.MakeRequest("POST", "crm.deal.add", map[string]any{"fields": dealEntityMap})
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Created
		deal = createdDeal
	}

	return ObjectStatus{
		CrmId:   deal["result"],
		Status:  status,
		Message: message,
	}, nil
}

func buildSearchFields(searchFilters map[string]any) string {
	stringFields, err := json.Marshal(searchFilters)
	if err != nil {
		fmt.Println("Error stringifying map:", err)
		return ""
	}
	return string(stringFields)
}

func processBitrixLead(client *BitrixClient, lead any, existingLead, rawData map[string]any, ownerId string) (*ObjectStatus, error) {
	exportedLead, exists := existingLead["lead"].(map[string]any)
	if exists && exportedLead["crm_id"] != nil {
		return createExistingStatus(exportedLead), nil
	}

	leadData, ok := lead.(map[string]any)
	if !ok {
		return nil, errors.New("invalid company data")
	}

	sentLead, err := sendBitrixLead(client, leadData, ownerId)
	if err != nil {
		return nil, err
	}

	if drivaID, exists := rawData["company_contact_id"].(string); exists {
		sentLead.DrivaContactId = drivaID
	}

	return &sentLead, nil
}

func sendBitrixLead(client *BitrixClient, mappedLeadData map[string]any, ownerId string) (ObjectStatus, error) {
	leadEntity, exists := mappedLeadData["entity"]
	if !exists {
		return ObjectStatus{
			Status:  Failed,
			Message: "lead entity not found in mapped lead data",
		}, errors.New("lead entity not found in mapped lead data")
	}
	leadEntityMap, isMap := leadEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{
			Status:  Failed,
			Message: "lead entity is not a map",
		}, errors.New("lead entity is not a map")
	}

	leadEntityMap["ASSIGNED_BY_ID"] = ownerId

	searchFilters := map[string]any{"TITLE": leadEntityMap["TITLE"]}
	existingLead, err := searchForExistingBitrixObject(client, "lead", searchFilters)
	if err != nil {
		return ObjectStatus{
			Status:  Failed,
			Message: err.Error(),
		}, err
	}

	searchFields := buildSearchFields(searchFilters)

	var lead map[string]any
	var status Status
	var message string
	if existingLead != nil {
		updatedLead, err := client.MakeRequest("POST", "crm.lead.update", map[string]any{"id": existingLead.(map[string]any)["ID"], "fields": leadEntityMap})
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		updatedLead["result"] = existingLead.(map[string]any)["ID"]
		status = Updated
		lead = updatedLead
		message = fmt.Sprintf("Searched fields: %s", searchFields)
	} else {
		createdLead, err := client.MakeRequest("POST", "crm.lead.add", map[string]any{"fields": leadEntityMap})
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Created
		lead = createdLead
	}

	return ObjectStatus{
		CrmId:   lead["result"],
		Status:  status,
		Message: message,
	}, nil
}

func processBitrixContact(client *BitrixClient, contact any, existingLead, rawData map[string]any, ownerId string) (*[]ObjectStatus, error) {
	exportedContact, exists := existingLead["contact"].(map[string]any)
	if exists && exportedContact["crm_id"] != nil {
		return &[]ObjectStatus{*createExistingStatus(exportedContact)}, nil
	}

	contactData, ok := contact.(map[string]any)
	if !ok {
		return nil, errors.New("invalid contact data")
	}

	sentContact, err := sendBitrixContact(client, contactData, ownerId)
	if err != nil {
		return nil, err
	}

	if drivaID, exists := rawData["profile_contact_id"].(string); exists {
		sentContact.DrivaContactId = drivaID
	}

	return &[]ObjectStatus{sentContact}, nil
}

func sendBitrixContact(client *BitrixClient, mappedContactData map[string]any, ownerId string) (ObjectStatus, error) {
	contactEntity, exists := mappedContactData["entity"]
	if !exists {
		return ObjectStatus{
			Status:  Failed,
			Message: "contact entity not found in mapped contact data",
		}, errors.New("contact entity not found in mapped contact data")
	}
	contactEntityMap, isMap := contactEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{
			Status:  Failed,
			Message: "contact entity is not a map",
		}, errors.New("contact entity is not a map")
	}

	contactEntityMap["ASSIGNED_BY_ID"] = ownerId

	var existingContact any
	var err error
	searchFilters := map[string]any{}
	var searchedEmails []string
	if emailValues, ok := contactEntityMap["EMAIL"].([]any); ok {
		for _, email := range emailValues {
			if value, ok := email.(map[string]any)["VALUE"]; ok && value != "" {
				searchFilters["EMAIL"] = value.(string)
				existingContact, err = searchForExistingBitrixObject(client, "contact", searchFilters)
				if err != nil {
					return ObjectStatus{
                        Status:  Failed,
                        Message: err.Error(),
                    }, err
                }
				searchedEmails = append(searchedEmails, value.(string))
				if existingContact != nil {
					break
				}
			}
		}
	}
	searchFilters["EMAIL"] = searchedEmails

	searchFields := buildSearchFields(searchFilters)

	var contact map[string]any
	var status Status
	var message string
	if existingContact != nil {
		updatedContact, err := client.MakeRequest("POST", "crm.contact.update", map[string]any{"id": existingContact.(map[string]any)["ID"], "fields": contactEntityMap})
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		updatedContact["result"] = existingContact.(map[string]any)["ID"]
		status = Updated
		contact = updatedContact
		message = fmt.Sprintf("Searched fields: %s", searchFields)
	} else {
		createdContact, err := client.MakeRequest("POST", "crm.contact.add", map[string]any{"fields": contactEntityMap})
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Created
		contact = createdContact
	}

	return ObjectStatus{
		CrmId:   contact["result"],
		Status:  status,
		Message: message,
	}, nil
}

func processBitrixContacts(client *BitrixClient, contacts any, existingLead, rawData map[string]any, ownerId string) (*[]ObjectStatus, error) {
	contactsData, ok := contacts.([]any)
	if !ok {
		return nil, errors.New("invalid contacts data")
	}

	var statuses []ObjectStatus
	for key, contact := range contactsData {
		var contactRawData map[string]any
		if key < len(rawData["profiles"].([]any)) {
			contactRawData = rawData["profiles"].([]any)[key].(map[string]any)
		}

		contactMap, ok := contact.(map[string]any)
		if !ok {
			continue
		}

		exportedContacts, exists := existingLead["contacts"].([]any)
		if exists {
			hasBeenSent := false
			for _, exportedContact := range exportedContacts {
				if exportedContact.(map[string]any)["driva_contact_id"] == contactRawData["profile_contact_id"] {
					statuses = append(statuses, *createExistingStatus(exportedContact.(map[string]any)))
					hasBeenSent = true
					break
				}
			}
			if hasBeenSent {
				continue
			}
		}

		sentContact, err := sendBitrixContact(client, contactMap, ownerId)
		if err != nil {
			return nil, err
		}

		if drivaID, exists := contactRawData["profile_contact_id"].(string); exists {
			sentContact.DrivaContactId = drivaID
		}
		statuses = append(statuses, sentContact)
	}

	return &statuses, nil
}
