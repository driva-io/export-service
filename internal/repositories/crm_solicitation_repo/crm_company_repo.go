package crm_solicitation_repo

import (
	"context"
	"encoding/json"
	"errors"
	"export-service/internal/repositories"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PgCrmSolicitationRepository struct {
	conn   *pgxpool.Pool
	logger *zap.Logger
}

func NewPgCrmSolicitationRepository(conn *pgxpool.Pool, logger *zap.Logger) *PgCrmSolicitationRepository {
	return &PgCrmSolicitationRepository{
		conn:   conn,
		logger: logger.Named("PgCrmSolicitationRepository"),
	}
}

func (r *PgCrmSolicitationRepository) GetById(ctx context.Context, id string) (Solicitation, error) {
	defer r.logger.Sync()

	rows, _ := r.conn.Query(ctx, getQuery, id)

	solicitation, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Solicitation])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Solicitation{}, repositories.NewSolicitationNotFoundError()
		} else {
			r.logger.Error("Got error when collecting one row", zap.Error(err), zap.Any("params", id))
		}

		return Solicitation{}, err
	}
	return solicitation, nil
}

func (r *PgCrmSolicitationRepository) Update(ctx context.Context, params UpdateExportedCompaniesParms, solicitationId string) (Solicitation, error) {
	defer r.logger.Sync()

	exportedCompanyBytes, err := json.Marshal(params.NewExportedCompany)
	if err != nil {
		return Solicitation{}, err
	}

	stringCnpj := fmt.Sprintf("%v", int(params.Cnpj.(float64)))

	rows, _ := r.conn.Query(ctx, updateExportedCompanies, stringCnpj, string(exportedCompanyBytes), solicitationId)

	solicitation, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Solicitation])
	if err != nil {
		r.logger.Error("Got error when collecting one row", zap.Error(err), zap.Any("params", params))
		if errors.Is(err, pgx.ErrNoRows) {
			return Solicitation{}, repositories.NewSolicitationNotFoundError()
		}

		return Solicitation{}, err
	}

	return solicitation, nil
}

func (r *PgCrmSolicitationRepository) UpdateStatus(ctx context.Context, newStatus SolicitationStatus, solicitationId string) (Solicitation, error) {
	defer r.logger.Sync()

	rows, _ := r.conn.Query(ctx, updateStatusQuery, newStatus, solicitationId)

	solicitation, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Solicitation])
	if err != nil {
		r.logger.Error("Got error when collecting one row", zap.Error(err), zap.Any("params", newStatus))
		if errors.Is(err, pgx.ErrNoRows) {
			return Solicitation{}, repositories.NewSolicitationNotFoundError()
		}

		return Solicitation{}, err
	}

	return solicitation, nil
}

func (r *PgCrmSolicitationRepository) IncrementCurrent(ctx context.Context, solicitationId string) (Solicitation, error) {
	defer r.logger.Sync()

	rows, _ := r.conn.Query(ctx, incrementCurrentQuery, solicitationId)

	solicitation, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Solicitation])
	if err != nil {
		r.logger.Error("Got error when collecting one row", zap.Error(err), zap.Any("params", solicitationId))
		if errors.Is(err, pgx.ErrNoRows) {
			return Solicitation{}, repositories.NewSolicitationNotFoundError()
		}

		return Solicitation{}, err
	}

	return solicitation, nil
}

func (r *PgCrmSolicitationRepository) Create(ctx context.Context, solicitation CreateSolicitation) (Solicitation, error) {
	defer r.logger.Sync()

	rows, err := r.conn.Query(ctx, createSolicitation, solicitation.ListId, solicitation.UserEmail, solicitation.OwnerId, solicitation.StageId, solicitation.PipelineId, solicitation.OverwriteData, solicitation.CreateDeal, solicitation.Current, solicitation.Total)
	if err != nil {
		return Solicitation{}, err
	}

	createdSolicitation, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Solicitation])
	if err != nil {
		r.logger.Error("Got error when collecting one row", zap.Error(err), zap.Any("params", solicitation))
		if errors.Is(err, pgx.ErrNoRows) {
			return Solicitation{}, repositories.NewSolicitationNotFoundError()
		}

		return Solicitation{}, err
	}

	return createdSolicitation, nil
}
