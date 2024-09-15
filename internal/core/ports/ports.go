package ports

import (
	"context"
	"export-service/internal/core/domain"
)

type PresentationSpecRepository interface {
	Get(ctx context.Context, userEmail string, userCompany string, service string, dataSource string) (domain.PresentationSpec, error)
}
