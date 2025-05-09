package usecases

import (
	"errors"
	"export-service/internal/adapters"
	"export-service/internal/repositories/presentation_spec_repo"
	"export-service/internal/writers"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSheetExportUseCase_downloadData(t *testing.T) {
	godotenv.Load("../../.env")
	config, err := pgxpool.ParseConfig(getPostgresConnStr())
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	conn, err := pgxpool.NewWithConfig(t.Context(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer conn.Close()

	p := presentation_spec_repo.NewPgPresentationSpecRepository(conn, zap.NewExample())
	u := getS3Uploader(zap.NewExample())
	m := adapters.NewDrivaMailer(zap.NewExample())
	s := SheetExportUseCase{presentationSpecRepo: p, uploader: u, dataWriter: &writers.ExcelWriter{}, downloader: &adapters.HTTPDownloader{}, logger: zap.NewExample(), mailer: m}
	t.Run("Should download data", func(t *testing.T) {
		r := ExportRequest{
			DataDownloadURL: getTestURL(t, serveJSON),
		}

		data, err := s.downloadData(r)
		require.NoError(t, err)

		assert.Lenf(t, data, 1, "Should have 1 company")

		assert.Equalf(t, "driva-tech", data[0]["public_id"], "Public ID should be driva-tech")
		assert.Lenf(t, data[0]["profiles"].([]any), 10, "Should have 10 profiles")
	})

	t.Run("Should fail if value is not an array", func(t *testing.T) {
		r := ExportRequest{
			DataDownloadURL: getTestURL(t, serveHTML),
		}

		_, err := s.downloadData(r)
		assert.Errorf(t, err, "Should fail if value is not an array")
	})

	t.Run("Should fail if status is not 200", func(t *testing.T) {
		r := ExportRequest{
			DataDownloadURL: getTestURL(t, serveErrorStatus),
		}

		_, err := s.downloadData(r)
		assert.Errorf(t, err, "Should fail if status is not 200")
	})

	t.Run("Should send sheet to email", func(t *testing.T) {
		r := ExportRequest{
			DataDownloadURL: "https://applications.s3.bhs.io.cloud.ovh.net/exports/request/e0775745-f311-4431-9db1-50ae452f4adf-e3c88c29",
			ListID:          "e0775745-f311-4431-9db1-50ae452f4adf",
			ListName:        "skate",
			DataSource:      "empresas",
			UserCompany:     "ramalho_company",
			UserEmail:       "henrique.ramalho@driva.com.br",
			UserName:        "Henrique",
		}
		_, err := s.Execute(r)
		require.NoError(t, err)
	})
}

func serveJSON(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("test_data/export_request.json")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<!DOCTYPE html><html><body></body></html>"))
}

func serveErrorStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("[]"))
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

func getS3Uploader(logger *zap.Logger) *adapters.S3Uploader {
	bucket := os.Getenv("S3_BUCKET")
	endpoint := os.Getenv("S3_ENDPOINT")
	folder := "exports/sheet"
	key := os.Getenv("S3_KEY")
	region := os.Getenv("S3_REGION")
	secretKey := os.Getenv("S3_SECRET_KEY")

	return adapters.NewS3Uploader(key, secretKey, endpoint, region, bucket, folder, logger)
}

func getTestURL(t *testing.T, handler http.HandlerFunc) string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Println("Error starting test server:", err)
		t.Fail()
	}

	port := listener.Addr().(*net.TCPAddr).Port
	server := &http.Server{Handler: handler}

	go func() {
		if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			log.Println("Test server closed with error:", err)
			t.Fail()
		}
	}()

	return fmt.Sprintf("http://0.0.0.0:%d", port)
}
