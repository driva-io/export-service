package crm_company_repo_test

import (
	"context"
	"database/sql"
	"export-service/internal/core/ports"
	"export-service/internal/repositories/crm_company_repo"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

func TestGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Setup postgres container

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:15-alpine",
		postgres.WithDatabase("exports"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts("seed.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
	url, _ := postgresContainer.ConnectionString(ctx)
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer conn.Close()
	logger, _ := zap.NewProduction()
	repo := crm_company_repo.NewPgCrmCompanyRepository(conn, logger)
	// Begin testing
	t.Run("Should return company", func(t *testing.T) {

		company := crm_company_repo.Company{
			Id:                 "123e4567-e89b-12d3-a456-426655440003",
			Crm:                "hubspot",
			Name:               sql.NullString{String: "Driva Teste F", Valid: true},
			RefreshToken:       sql.NullString{String: "refresh_token_1", Valid: true},
			AccessToken:        sql.NullString{String: "access_token_1", Valid: true},
			ExpiresIn:          sql.NullString{String: "3600", Valid: true},
			RefreshedAt:        sql.NullString{String: "2022-01-01 00:00:00 -03:00", Valid: true},
			Environment:        sql.NullString{String: "production", Valid: true},
			Token:              sql.NullString{String: "token_1", Valid: true},
			Webhook:            sql.NullString{String: "https://webhook.url/1", Valid: true},
			Email:              sql.NullString{String: "francisco.becheli@driva.com.br", Valid: true},
			Password:           sql.NullString{String: "", Valid: false},
			InstanceUrl:        sql.NullString{String: "https://instance.url/1", Valid: true},
			Merge:              sql.NullString{String: "merge_A", Valid: true},
			Mapping:            sql.NullString{String: "", Valid: false},
			MappingLinkedin:    sql.NullString{String: "", Valid: false},
			CompanyId:          sql.NullString{String: "company_id_1", Valid: true},
			UserWhoInstalledId: sql.NullString{String: "user_id_1", Valid: true},
			CreatedAt:          time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt:          time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			WorkspaceId:        sql.NullString{String: "workspace_1", Valid: true},
		}

		result, err := repo.Get(ctx, ports.CrmCompanyQueryParams{Crm: "hubspot", WorkspaceId: "workspace_1"})

		require.NoError(t, err)
		require.Equal(t, company, result)
	})

	// t.Run("Should return error if no company found", func(t *testing.T) {

	// 	_, err := repo.Get(ctx, ports.CrmCompanyQueryParams{Crm: "wrong_crm", Company: "Wrong Company"})

	// 	require.Error(t, err)
	// })

	// t.Run("Should return error if invalid params", func(t *testing.T) {

	// 	_, err := repo.Get(ctx, ports.CrmCompanyQueryParams{})

	// 	var invalidErr ports.InvalidQueryParamsError
	// 	require.Error(t, err)
	// 	require.ErrorAs(t, err, &invalidErr)
	// })
}
