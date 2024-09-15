package presentation_spec_repo

import (
	"context"
	"export-service/internal/core/domain"
	"export-service/internal/core/ports"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PgPresentationSpecRepository struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func NewPgPresentationSpecRepository(conn *pgx.Conn, logger *zap.Logger) *PgPresentationSpecRepository {
	return &PgPresentationSpecRepository{
		conn:   conn,
		logger: logger.Named("PgPresentationSpecRepository"),
	}
}

func (r *PgPresentationSpecRepository) Get(ctx context.Context, params ports.PresentationSpecQueryParams) (domain.PresentationSpec, error) {
	defer r.logger.Sync()

	if params.UserEmail == "" || params.UserCompany == "" || params.Service == "" || params.DataSource == "" {
		return domain.PresentationSpec{}, ports.ErrInvalidParams
	}

	rows, _ := r.conn.Query(ctx, getQuery, params.UserEmail, params.UserCompany, params.Service, params.DataSource)

	spec, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.PresentationSpec])
	if err != nil {
		r.logger.Error("Got error when collecting one row", zap.Error(err), zap.Any("params", params))
		return domain.PresentationSpec{}, err
	}
	return spec, nil
}
