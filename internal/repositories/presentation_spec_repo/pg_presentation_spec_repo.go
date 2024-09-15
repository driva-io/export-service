package presentation_spec_repo

import (
	"context"
	"export-service/internal/core/domain"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PgPresentationSpecRepository struct {
	conn   *pgx.Conn
	logger *zap.SugaredLogger
}

func NewPgPresentationSpecRepository(conn *pgx.Conn) *PgPresentationSpecRepository {
	logger := zap.NewExample().Sugar().Named("PgPresentationSpecRepository")
	return &PgPresentationSpecRepository{
		conn:   conn,
		logger: logger,
	}
}

func (r *PgPresentationSpecRepository) Get(ctx context.Context, userEmail string, userCompany string, service string, dataSource string) (domain.PresentationSpec, error) {
	defer r.logger.Sync()
	rows, _ := r.conn.Query(ctx, getQuery, userEmail, userCompany, service, dataSource)

	spec, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.PresentationSpec])
	if err != nil {
		r.logger.Errorw("Got error when collecting one row", "error", err, "userEmail", userEmail, "userCompany", userCompany, "service", service, "dataSource", dataSource)
		return domain.PresentationSpec{}, err
	}
	return spec, nil
}
