package main

import (
	"context"
	"export-service/api/routes"
	"export-service/internal/gateways"
	"export-service/internal/repositories/crm_company_repo"
	"export-service/internal/repositories/presentation_spec_repo"
	srv "export-service/internal/server"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	server := srv.New()

	ctx := context.Background()

	logger := zap.NewExample()
	conn, err := pgx.Connect(ctx, getPostgresConnStr())

	auth := &gateways.HTTPAuthService{HttpClient: &srv.NetHttpClient{}}

	presentationSpecRepo := presentation_spec_repo.NewPgPresentationSpecRepository(conn, logger)
	crmCompanyRepo := crm_company_repo.NewPgCrmCompanyRepository(conn, logger)

	routes.RegisterServerRoutes(server, auth)
	routes.RegisterPresentationSpecRoutes(server, presentationSpecRepo, auth)
	routes.RegisterCrmRoutes(server, auth, crmCompanyRepo, presentationSpecRepo)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err = server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
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
