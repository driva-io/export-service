package ports

import (
	"context"
	"export-service/internal/core/domain"
)

type CrmCompanyQueryParams struct {
	Crm         string
	WorkspaceId string
}

type CrmGetByCompanyNameQueryParams struct {
	Crm         string
	CompanyName string
}

type CrmAddHubspotCompanyQueryParams struct {
	WorkspaceId  string
	UserId       string
	RefreshToken string
	AccessToken  string
	ExpiresIn string
}

type PresentationSpecQueryParams struct {
	UserEmail   string
	UserCompany string
	Service     string
	DataSource  string
}

type PresentationSpecAddBody struct {
	PresentationSpec domain.PresentationSpecSpec           `json:"spec" validate:"required"`
	SpecOptions      []domain.PresentationSpecSheetOptions `json:"sheet_options" validate:"required"`
}

type PresentationSpecPatchKey struct {
	PresentationSpec map[string]any                           `json:"spec" validate:"required"`
	SpecOptions      domain.PresentationSpecPatchSheetOptions `json:"sheet_options" validate:"required"`
}

type PresentationSpecRepository interface {
	Get(ctx context.Context, params PresentationSpecQueryParams) (domain.PresentationSpec, error)
	Add(ctx context.Context, params PresentationSpecQueryParams, body PresentationSpecAddBody) (domain.PresentationSpec, error)
}

type DataWriter interface {
	Write(data []map[string]any, spec domain.PresentationSpec) (string, error)
}

type Downloader interface {
	Download(url string) ([]byte, error)
}

type Uploader interface {
	Upload(fileName, path string) (string, error)
}

type Mailer interface {
	SendEmail(userEmail, userName, templateId, link string) error
}
