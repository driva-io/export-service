package usecases

import (
	"context"
	"encoding/json"
	"export-service/internal/core/domain"
	"export-service/internal/core/ports"
	"export-service/internal/services/data_presenter"
	"fmt"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

type SheetExportUseCase struct {
	dataWriter           ports.DataWriter
	downloader           ports.Downloader
	uploader             ports.Uploader
	presentationSpecRepo ports.PresentationSpecRepository
	mailer               ports.Mailer
	logger               *zap.Logger
}

func NewSheetExportUseCase(dataWriter ports.DataWriter, downloader ports.Downloader, uploader ports.Uploader, presentationSpecRepo ports.PresentationSpecRepository, mailer ports.Mailer, logger *zap.Logger) *SheetExportUseCase {
	return &SheetExportUseCase{
		dataWriter:           dataWriter,
		downloader:           downloader,
		uploader:             uploader,
		presentationSpecRepo: presentationSpecRepo,
		mailer:               mailer,
		logger:               logger,
	}
}

func (s *SheetExportUseCase) Execute(request ExportRequest) (string, error) {
	data, err := s.downloadData(request)
	if err != nil {
		s.logError("Error when downloading data", err, request)
		return "", err
	}

	spec, err := s.getPresentationSpec(request)
	if err != nil {
		s.logError("Error when getting presentation spec", err, request)
		return "", err
	}

	presentedData, err := s.applyPresentationSpec(request, data, spec)
	if err != nil {
		s.logError("Error when applying presentation spec", err, request)
		return "", err
	}

	path, err := s.writeSheet(request, presentedData, spec)
	if err != nil {
		s.logError("Error when writing data", err, request)
		return "", err
	}

	url, err := s.uploadSheet(request, path)
	if err != nil {
		s.logError("Error when uploading sheet", err, request)
		return "", err
	}

	err = s.sendEmail(request, url)
	if err != nil {
		s.logError("Error when sending email", err, request)
		return "", err
	}

	return url, nil
}

func (s *SheetExportUseCase) downloadData(request ExportRequest) ([]map[string]any, error) {
	s.logInfo("Downloading data", request)
	bytes, err := s.downloader.Download(request.DataDownloadURL)
	if err != nil {
		return nil, err
	}

	var data []map[string]any
	err = json.Unmarshal(bytes, &data)
	return data, err
}

func (s *SheetExportUseCase) getPresentationSpec(request ExportRequest) (domain.PresentationSpec, error) {
	s.logInfo("Getting presentation spec", request)

	return s.presentationSpecRepo.Get(context.Background(), ports.PresentationSpecQueryParams{
		UserEmail:   request.UserEmail,
		UserCompany: request.UserCompany,
		Service:     "sheet",
		DataSource:  request.DataSource,
	})
}

func (s *SheetExportUseCase) applyPresentationSpec(request ExportRequest, data []map[string]any, spec domain.PresentationSpec) ([]map[string]any, error) {
	s.logInfo("Applying presentation spec", request)
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

func (s *SheetExportUseCase) writeSheet(request ExportRequest, data []map[string]any, spec domain.PresentationSpec) (string, error) {
	s.logInfo("Writing sheet", request)
	return s.dataWriter.Write(data, spec)
}

func (s *SheetExportUseCase) uploadSheet(request ExportRequest, path string) (string, error) {
	s.logInfo("Uploading sheet", request)
	return s.uploader.Upload(fmt.Sprintf("(DRIVA %s) %s.xlsx", uuid.NewString()[:8], request.ListName), path)
}

func (s *SheetExportUseCase) sendEmail(request ExportRequest, url string) error {
	s.logInfo("Sending email", request)
	return s.mailer.SendEmail(request.UserEmail, request.UserName, "d-b1c0014b01eb410a8c4b8112e4418a3f", url)
}

func (s *SheetExportUseCase) logInfo(message string, request ExportRequest) {
	s.logger.Info(message, zap.Any("request", request))
}

func (s *SheetExportUseCase) logError(message string, err error, request ExportRequest) {
	s.logger.Error(message, zap.Error(err), zap.Any("request", request))
}
