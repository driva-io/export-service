package crm_company_repo

import (
	"context"
	"errors"
	"export-service/internal/core/ports"
	"export-service/internal/repositories"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PgCrmCompanyRepository struct {
	conn   *pgx.Conn
	logger *zap.Logger
}

func NewPgCrmCompanyRepository(conn *pgx.Conn, logger *zap.Logger) *PgCrmCompanyRepository {
	return &PgCrmCompanyRepository{
		conn:   conn,
		logger: logger.Named("PgCrmCompanyRepository"),
	}
}

func (r *PgCrmCompanyRepository) Get(ctx context.Context, params ports.CrmCompanyQueryParams) (Company, error) {
	defer r.logger.Sync()

	if params.Company == "" || params.Crm == "" {
		return Company{}, ports.NewInvalidQueryParamsError()
	}

	rows, err := r.conn.Query(ctx, getQuery, params.Crm, params.Company)
	if err != nil {
		r.logger.Error("Failed to execute query", zap.Error(err), zap.Any("params", params))
		return Company{}, err
	}
	defer rows.Close()

	company, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Company])
	if err != nil {
		r.logger.Error("Got error when collecting one row", zap.Error(err), zap.Any("params", params))
		if errors.Is(err, pgx.ErrNoRows) {
			return Company{}, repositories.NewCompanyNotFoundError()
		} else if errors.Is(err, pgx.ErrTooManyRows) {
			return Company{}, repositories.NewCompanyNotUniqueError()
		}

		return Company{}, err
	}
	return company, nil
}

func (r *PgCrmCompanyRepository) AddHubspot(ctx context.Context, params ports.CrmAddHubspotCompanyQueryParams) (Company, error) {
	defer r.logger.Sync()

	if params.Company == "" || params.RefreshToken == "" || params.AccessToken == "" || params.UserId == "" || params.WorkspaceId == "" {
		return Company{}, ports.NewInvalidQueryParamsError()
	}

	rows, err := r.conn.Query(ctx, addHubspotQuery, params.Company, params.UserId, params.WorkspaceId, params.RefreshToken, params.AccessToken)
	if err != nil {
		r.logger.Error("Failed to execute query", zap.Error(err), zap.Any("params", params))
		return Company{}, err
	}
	defer rows.Close()

	company, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Company])
	if err != nil {
		return Company{}, err
	}
	return company, nil
}

// func (r *PgPresentationSpecRepository) GetById(ctx context.Context, id string) (domain.PresentationSpec, error) {
// 	defer r.logger.Sync()

// 	rows, _ := r.conn.Query(ctx, getByIdQuery, id)

// 	spec, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.PresentationSpec])
// 	if err != nil {
// 		r.logger.Error("Got error when collecting one row", zap.Error(err), zap.Any("params", id))
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return domain.PresentationSpec{}, repositories.NewPresentationSpecNotFoundError()
// 		} else if errors.Is(err, pgx.ErrTooManyRows) {
// 			return domain.PresentationSpec{}, repositories.NewPresentationSpecNotUniqueError()
// 		}

// 		return domain.PresentationSpec{}, err
// 	}
// 	return spec, nil
// }

// func (r *PgPresentationSpecRepository) Patch(ctx context.Context, id string, body ports.PresentationSpecAddBody) (domain.PresentationSpec, error) {
// 	defer r.logger.Sync()

// 	presentationSpec := body.PresentationSpec
// 	sheetOptions := body.SpecOptions

// 	if id == "" {
// 		return domain.PresentationSpec{}, ports.NewInvalidQueryParamsError()
// 	}

// 	for _, value := range sheetOptions {
// 		keyName := value.Key
// 		_, exists := presentationSpec[keyName]
// 		if !exists {
// 			return domain.PresentationSpec{}, errors.New("key " + keyName + "in sheet_options is not present in spec.")
// 		}
// 	}

// 	for key := range presentationSpec {
// 		found := false
// 		for _, value := range sheetOptions {
// 			if key == value.Key {
// 				if found {
// 					return domain.PresentationSpec{}, errors.New("duplicate key " + key + " in sheet_options.")
// 				}
// 				found = true
// 			}
// 		}
// 		if !found {
// 			return domain.PresentationSpec{}, errors.New("sheet " + key + "in spec is not present in sheet_options.")
// 		}
// 	}

// 	tx, err := r.conn.Begin(ctx)
// 	if err != nil {
// 		r.logger.Error("Failed to start transaction", zap.Error(err))
// 		return domain.PresentationSpec{}, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
// 				r.logger.Error("Failed to roll back transaction", zap.Error(rollbackErr))
// 			}
// 		} else {
// 			if commitErr := tx.Commit(ctx); commitErr != nil {
// 				r.logger.Error("Failed to commit transaction", zap.Error(commitErr))
// 			}
// 		}
// 	}()

// 	if _, err := tx.Exec(ctx, deleteSpecsQuery, id); err != nil {
// 		r.logger.Error("Got error when deleting specs", zap.Error(err), zap.Any("params", id))
// 		return domain.PresentationSpec{}, err
// 	}

// 	if _, err := tx.Exec(ctx, deleteSheetOptionsQuery, id); err != nil {
// 		r.logger.Error("Got error when deleting sheet options", zap.Error(err), zap.Any("params", id))
// 		return domain.PresentationSpec{}, err
// 	}

// 	if _, err := tx.Exec(ctx, patchBasicInfo, id); err != nil {
// 		r.logger.Error("Got error patching basic_info", zap.Error(err), zap.Any("params", id))
// 		return domain.PresentationSpec{}, err
// 	}

// 	for tab, tabSpec := range presentationSpec {
// 		var correspondingOptions domain.PresentationSpecSheetOptions
// 		for _, options := range sheetOptions {
// 			if tab == options.Key {
// 				correspondingOptions = options
// 				break
// 			}
// 		}

// 		if _, err := tx.Exec(ctx, addOptionsQuery, id, correspondingOptions.Key, correspondingOptions.ActiveColumns, correspondingOptions.Position, correspondingOptions.ShouldExplode); err != nil {
// 			r.logger.Error("Got error when inserting options", zap.Error(err), zap.Any("params", id))
// 			return domain.PresentationSpec{}, err
// 		}

// 		if _, err := tx.Exec(ctx, addSpecQuery, id, correspondingOptions.Key, tabSpec); err != nil {
// 			r.logger.Error("Got error when inserting specs", zap.Error(err), zap.Any("params", id))
// 			return domain.PresentationSpec{}, err
// 		}
// 	}

// 	patchedSpec, _ := r.GetById(ctx, id)

// 	return patchedSpec, nil
// }

// func (r *PgPresentationSpecRepository) PatchKey(ctx context.Context, id string, key string, body ports.PresentationSpecPatchKey) (domain.PresentationSpec, error) {
// 	defer r.logger.Sync()

// 	if id == "" || key == "" {
// 		return domain.PresentationSpec{}, ports.NewInvalidQueryParamsError()
// 	}

// 	spec := body.PresentationSpec
// 	options := body.SpecOptions

// 	rows, err := r.conn.Query(ctx, patchKeyOptions, options.Key, options.ActiveColumns, options.Position, options.ShouldExplode, id, key)
// 	if err != nil {
// 		r.logger.Error("Got error when updating sheet options", zap.Error(err), zap.Any("params", id))
// 		return domain.PresentationSpec{}, err
// 	}
// 	rows.Close()

// 	rows, err = r.conn.Query(ctx, patchKeySpec, spec, options.Key, id, key)
// 	if err != nil {
// 		r.logger.Error("Got error when updating key spec", zap.Error(err), zap.Any("params", id))
// 		return domain.PresentationSpec{}, err
// 	}
// 	rows.Close()

// 	updatedSpec, _ := r.GetById(ctx, id)

// 	return updatedSpec, nil
// }

// func (r *PgPresentationSpecRepository) Delete(ctx context.Context, id string) error {
// 	defer r.logger.Sync()

// 	if id == "" {
// 		return ports.NewInvalidQueryParamsError()
// 	}

// 	rows, err := r.conn.Query(ctx, deleteQuery, id)
// 	if err != nil {
// 		r.logger.Error("Got error when deleting basic info", zap.Error(err), zap.Any("params", id))
// 		return err
// 	}
// 	rows.Close()

// 	return nil
// }
