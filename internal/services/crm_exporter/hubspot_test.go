package crm_exporter

import (
	"context"
	"export-service/internal/repositories/crm_company_repo"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSendLead(t *testing.T) {
	t.Parallel()

	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()

	lead := map[string]any{
		"company": map[string]any{
			"entity": map[string]any{
				"name":         "empresa teste",
				"cnpj_empresa": 124,
				"state":        "PR",
			},
		},
		"deal": map[string]any{
			"entity": map[string]any{
				"dealname": "deal teste",
			},
		},
		"contacts": []map[string]any{
			{
				"entity": map[string]any{
					"firstname": "contato teste 1",
					"email":     "fulano@teste.com",
				},
			},
			{
				"entity": map[string]any{
					"firstname": "contato teste 2",
					"email":     "fulano2@teste.com",
				},
			},
		},
	}

	logger := zap.NewExample()
	config, err := pgxpool.ParseConfig(getPostgresConnStr())
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer conn.Close()

	crmCompanyRepo := crm_company_repo.NewPgCrmCompanyRepository(conn, logger)
	hubspotService := NewHubspotService(crmCompanyRepo)
	client, _ := hubspotService.Authorize(ctx, "Driva Teste F")

	t.Run("Should send lead", func(t *testing.T) {
		t.Skip("Skipping hubspot prod test")
		configs := map[string]any{
			"owner_id":    "718806932",
			"pipeline_id": "default",
			"stage_id":    "appointmentscheduled",
			"create_deal": true,
		}
		result, err := hubspotService.SendLead(client, lead, map[string]any{}, configs, map[string]any{})

		println(result.Company)
		require.NoError(t, err)
	})
}

func getPostgresConnStr() string {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_EXPORTS_DATABASE")
	escapedUser := url.QueryEscape(user)
	escapedPassword := url.QueryEscape(password)

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", escapedUser, escapedPassword, host, port, dbname)
}
