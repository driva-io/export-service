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
