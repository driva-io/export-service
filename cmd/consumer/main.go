package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"net/url"

	"context"
	"encoding/json"
	"export-service/internal/adapters"
	"export-service/internal/messaging"
	"export-service/internal/repositories/presentation_spec_repo"
	"export-service/internal/usecases"
	"export-service/internal/writers"
	"github.com/jackc/pgx/v5"
	"log"
	"os"

	"go.uber.org/zap"
)

func main() {
	con, err := messaging.Connect(os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer con.Close()

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, getPostgresConnStr())
	failOnError(err, "Failed to connect to database")

	defer conn.Close(ctx)

	logger := zap.NewExample()
	client := messaging.NewRabbitMQClient(con, logger)
	defer client.Close()

	sheetUc := getSheetUseCase(logger, conn)

	// Create queues
	exports := "exports.excel"
	failOnError(client.CreateQueue("exports.excel", false), "Failed to create exports queue")
	failOnError(client.CreateQueue("exports.results.excel", false), "Failed to create exports result queue")

	exportsBus, err := client.Consume(exports)
	failOnError(err, "Failed to consume bus")

	var blocking chan struct{}

	go func() {
		for d := range exportsBus {
			logger.Info("Received message on exports queue", zap.Any("message", d))

			// Unmarshal Body into DTO
			var req usecases.ExportRequest
			if err = json.Unmarshal(d.Body, &req); err != nil {
				logger.Error("Failed to unmarshal message", zap.Error(err))
				failOnError(d.Nack(false, false), "Failed to nack message")
			}

			downloadUrl, err := sheetUc.Execute(req)
			publishResult(client, logger, req, downloadUrl, err)

			failOnError(d.Ack(false), "Failed to ack message")
		}
	}()

	logger.Info("Consuming messages, press CTRL+C to stop")
	// Blocks forever
	<-blocking
}

func publishResult(c *messaging.RabbitClient, logger *zap.Logger, req usecases.ExportRequest, downloadUrl string, err error) {
	response := struct {
		ListID      string `json:"list_id,omitempty"`
		DownloadUrl string `json:"download_url,omitempty"`
		Error       string `json:"error,omitempty"`
	}{
		ListID:      req.ListID,
		DownloadUrl: downloadUrl,
	}

	if err != nil {
		response.Error = err.Error()
		logger.Error("Failed to execute use case", zap.Error(err))
	} else {
		logger.Info("Successfully executed use case")
	}

	b, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response", zap.Error(err))
	}

	failOnError(c.Publish("exports.results.excel", b), "Failed to publish response message")
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

func getSheetUseCase(logger *zap.Logger, conn *pgx.Conn) *usecases.SheetExportUseCase {
	bucket := os.Getenv("S3_BUCKET")
	endpoint := os.Getenv("S3_ENDPOINT")
	folder := "exports/sheet"
	key := os.Getenv("S3_KEY")
	region := os.Getenv("S3_REGION")
	secretKey := os.Getenv("S3_SECRET_KEY")

	uploader := adapters.NewS3Uploader(key, secretKey, endpoint, region, bucket, folder, logger)
	mailer := adapters.NewDrivaMailer(logger)

	specRepo := presentation_spec_repo.NewPgPresentationSpecRepository(conn, logger)
	return usecases.NewSheetExportUseCase(&writers.ExcelWriter{}, &adapters.HTTPDownloader{}, uploader, specRepo, mailer, logger)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
