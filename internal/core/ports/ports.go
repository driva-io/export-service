package ports

import (
	"context"
	"errors"
	"export-service/internal/core/domain"
)

type PresentationSpecQueryParams struct {
	UserEmail   string
	UserCompany string
	Service     string
	DataSource  string
}

var ErrInvalidParams = errors.New("invalid query params provided")

type PresentationSpecRepository interface {
	Get(ctx context.Context, params PresentationSpecQueryParams) (domain.PresentationSpec, error)
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
