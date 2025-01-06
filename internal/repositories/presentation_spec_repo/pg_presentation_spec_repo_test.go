package presentation_spec_repo_test

import (
	"context"
	"export-service/internal/core/domain"
	"export-service/internal/core/ports"
	"export-service/internal/repositories/presentation_spec_repo"
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
	repo := presentation_spec_repo.NewPgPresentationSpecRepository(conn, logger)
	// Begin testing
	t.Run("Should return user's custom spec", func(t *testing.T) {

		presentationSpec := domain.PresentationSpec{
			ID:          "123e4567-e89b-12d3-a456-426655440000",
			UserEmail:   "victor@driva.com.br",
			UserCompany: "Driva",
			Base:        "empresas",
			Service:     "enrichment_test",
			Version:     2,
			CreatedAt:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			SheetOptions: []domain.PresentationSpecSheetOptions{
				{
					Key:           "RFB",
					ActiveColumns: []string{"CNPJ"},
					Position:      0,
					ShouldExplode: false,
				},
			},
			Spec: domain.PresentationSpecSpec{
				"RFB": map[string]any{"CNPJ": "cnpj"},
			},
		}

		result, err := repo.Get(ctx, ports.PresentationSpecQueryParams{UserEmail: "victor@driva.com.br", UserCompany: "Driva", Service: "enrichment_test", DataSource: "empresas"})

		require.NoError(t, err)
		require.Equal(t, presentationSpec, result)
	})

	t.Run("Should return the default spec for user without custom spec", func(t *testing.T) {

		defaultSheetOptions := []domain.PresentationSpecSheetOptions{
			{
				Key:           "RFB",
				ActiveColumns: []string{"CNPJ", "Nome"},
				Position:      0,
				ShouldExplode: false,
			},
		}
		defaultSpec := domain.PresentationSpecSpec{
			"RFB": map[string]any{"CNPJ": "cnpj", "Nome": "razao_social"},
		}

		result, err := repo.Get(ctx, ports.PresentationSpecQueryParams{UserEmail: "user_sem_spec@driva.com.br", UserCompany: "Driva", Service: "enrichment_test", DataSource: "empresas"})

		require.NoError(t, err)
		require.Equal(t, result.IsDefault, true)
		require.Equal(t, defaultSheetOptions, result.SheetOptions)
		require.Equal(t, defaultSpec, result.Spec)

	})

	t.Run("Should return error if not spec found", func(t *testing.T) {

		_, err := repo.Get(ctx, ports.PresentationSpecQueryParams{UserEmail: "user_sem_spec@driva.com.br", UserCompany: "Driva", Service: "enrichment_test", DataSource: "base_que_nao_existe"})

		require.Error(t, err)
	})

	t.Run("Should return error if invalid params", func(t *testing.T) {

		_, err := repo.Get(ctx, ports.PresentationSpecQueryParams{})

		var invalidErr ports.InvalidQueryParamsError
		require.Error(t, err)
		require.ErrorAs(t, err, &invalidErr)
	})
}

func TestAdd(t *testing.T) {
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
	repo := presentation_spec_repo.NewPgPresentationSpecRepository(conn, logger)

	sheetOptions := []domain.PresentationSpecSheetOptions{
		{
			Key:           "RFB",
			ActiveColumns: []string{"CNPJ", "Razão Social", "Endereço"},
			Position:      0,
			ShouldExplode: false,
		},
		{
			Key:           "Telefones",
			ActiveColumns: []string{"CNPJ", "Razão Social", "Telefone Completo"},
			Position:      1,
			ShouldExplode: true,
		},
	}

	spec := domain.PresentationSpecSpec{
		"RFB":       {"CNPJ": "cnpj", "Razão Social": "razao_social", "Endereço": "endereco"},
		"Telefones": {"CNPJ": "cnpj", "Razão Social": "razao_social", "Telefone Completo": "1111-1111"},
	}

	// Begin testing
	t.Run("Should add user's custom spec", func(t *testing.T) {

		result, err := repo.Add(ctx, ports.PresentationSpecQueryParams{UserEmail: "francisco.becheli@driva.com.br", UserCompany: "Driva", Service: "enrichment", DataSource: "empresas"}, ports.PresentationSpecAddBody{SpecOptions: sheetOptions, PresentationSpec: spec})

		require.NoError(t, err)
		require.Equal(t, result.SheetOptions, sheetOptions)
		require.Equal(t, result.Spec, spec)
		require.Equal(t, result.UserEmail, "francisco.becheli@driva.com.br")
		require.Equal(t, result.UserCompany, "Driva")
		require.Equal(t, result.Service, "enrichment")
		require.Equal(t, result.Base, "empresas")
	})

	t.Run("Should return error if invalid params", func(t *testing.T) {

		_, err := repo.Add(ctx, ports.PresentationSpecQueryParams{UserCompany: "Driva", Service: "enrichment", DataSource: "empresas"}, ports.PresentationSpecAddBody{SpecOptions: sheetOptions, PresentationSpec: spec})

		var invalidErr ports.InvalidQueryParamsError
		require.Error(t, err)
		require.ErrorAs(t, err, &invalidErr)
	})
}

func TestDelete(t *testing.T) {
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
	repo := presentation_spec_repo.NewPgPresentationSpecRepository(conn, logger)

	sheetOptions := []domain.PresentationSpecSheetOptions{
		{
			Key:           "RFB",
			ActiveColumns: []string{"CNPJ", "Razão Social", "Endereço"},
			Position:      0,
			ShouldExplode: false,
		},
		{
			Key:           "Telefones",
			ActiveColumns: []string{"CNPJ", "Razão Social", "Telefone Completo"},
			Position:      1,
			ShouldExplode: true,
		},
	}

	spec := domain.PresentationSpecSpec{
		"RFB":       {"CNPJ": "cnpj", "Razão Social": "razao_social", "Endereço": "endereco"},
		"Telefones": {"CNPJ": "cnpj", "Razão Social": "razao_social", "Telefone Completo": "1111-1111"},
	}

	result, _ := repo.Add(ctx, ports.PresentationSpecQueryParams{UserEmail: "francisco.becheli@driva.com.br", UserCompany: "Driva", Service: "enrichment", DataSource: "empresas"}, ports.PresentationSpecAddBody{SpecOptions: sheetOptions, PresentationSpec: spec})

	specId := result.ID

	// Begin testing
	t.Run("Should delete user's custom spec", func(t *testing.T) {

		err := repo.Delete(ctx, specId)

		require.NoError(t, err)
	})
}
