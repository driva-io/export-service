package main

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/module/apmzap/v2"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap/zapcore"

	"context"
	"encoding/json"
	"export-service/internal/adapters"
	"export-service/internal/messaging"
	"export-service/internal/repositories/crm_company_repo"
	"export-service/internal/repositories/crm_solicitation_repo"
	"export-service/internal/repositories/presentation_spec_repo"
	"export-service/internal/server"
	"export-service/internal/usecases"
	"export-service/internal/writers"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"
)

var mainLogger = getLogger()

func main() {
	con, err := messaging.Connect(os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer con.Close()

	ctx := context.Background()

	config, err := pgxpool.ParseConfig(getPostgresConnStr())
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer conn.Close()

	client := messaging.NewRabbitMQClient(con, mainLogger)
	defer client.Close()

	// Create queues
	dlx := "exports_dead_letter_exchange"
	crmRKey := "exports.crm.dlq"
	exports := "exports.excel"
	failOnError(client.CreateQueue("exports.excel", nil, nil), "Failed to create exports queue")
	failOnError(client.CreateQueue("exports.results.excel", nil, nil), "Failed to create exports result queue")
	crm := "exports.crm"
	failOnError(client.CreateQueue(crm, &dlx, &crmRKey), "Failed to create exports.crm queue")

	go func() {
		for {
			exportsBus, err := client.Consume(exports)
			failOnError(err, "Failed to consume bus")

			for d := range exportsBus {
				// if conn.IsClosed() {
				// 	conn, err = pgx.Connect(ctx, getPostgresConnStr())
				// 	failOnError(err, "Failed to connect to database")
				// }

				handleExportRequest(d, conn, client)
			}

			mainLogger.Warn("Queue closed, retrying in 60 seconds")
			time.Sleep(60 * time.Second)
		}
	}()

	// go func() {
	// 	for {
	// 		crmBus, err := client.Consume(crm)
	// 		failOnError(err, "Failed to consume bus")

	// 		for d := range crmBus {
	// 			// if conn.IsClosed() {
	// 			// 	conn, err = pgx.Connect(ctx, getPostgresConnStr())
	// 			// 	failOnError(err, "Failed to connect to database")
	// 			// }

	// 			handleCrmExportRequest(d, conn)
	// 		}

	// 		mainLogger.Warn("Queue closed, retrying in 60 seconds")
	// 		time.Sleep(60 * time.Second)
	// 	}
	// }()

	mainLogger.Info("Consuming messages, press CTRL+C to stop")
	// Blocks forever
	<-make(chan struct{})
}

func handleCrmExportRequest(d amqp.Delivery, conn *pgxpool.Pool) {
	ctx := getMessageContext(d)
	defer func(ctx context.Context) {
		tx := apm.TransactionFromContext(ctx)
		if tx != nil {
			tx.End()
		}
	}(ctx)

	logger := mainLogger.With(apmzap.TraceContext(ctx)...)
	logger.Info("Received message on crm exports queue", zap.Any("message", d))

	var req usecases.CrmExportRequest
	if err := json.Unmarshal(d.Body, &req); err != nil {
		logger.Error("Failed to unmarshal message", zap.Error(err))
		failOnError(d.Nack(false, false), "Failed to nack message")
	}

	headers := d.Headers
	total := int64(0)
	if rawTotal, ok := headers["total"]; ok {
		switch v := rawTotal.(type) {
		case int:
			total = int64(v)
		case int64:
			total = v
		case float64:
			total = int64(v)
		default:
			logger.Warn("Unexpected type for total", zap.Any("value", rawTotal))
			failOnError(d.Nack(false, false), "Failed to nack message")
			return
		}
	} else {
		logger.Warn("Total is missing in headers")
		failOnError(d.Nack(false, false), "Failed to nack message")
		return
	}

	configs := map[string]any{
		"crm":            headers["crm"],
		"pipeline_id":    headers["pipeline_id"],
		"stage_id":       headers["stage_id"],
		"owner_id":       headers["owner_id"],
		"create_deal":    headers["create_deal"],
		"overwrite_data": headers["overwrite_data"],
		"total":          total,
		//Add other crm configs
	}

	CrmUc := getCrmUseCase(logger, conn)
	err := CrmUc.Execute(req, configs)
	if err != nil {
		retriable := CrmUc.IsRetriable(err)
		logger.Error("Error executing CRM request",
			zap.Error(err),
			zap.Bool("retriable", retriable),
		)
		failOnError(d.Nack(false, false), "Failed to nack message")
	} else {
		failOnError(d.Ack(false), "Failed to ack message")
	}
}

func getCrmUseCase(logger *zap.Logger, conn *pgxpool.Pool) *usecases.CrmExportUseCase {
	mailer := adapters.NewDrivaMailer(logger)
	specRepo := presentation_spec_repo.NewPgPresentationSpecRepository(conn, logger)
	companyRepo := crm_company_repo.NewPgCrmCompanyRepository(conn, logger)
	solicitationRepo := crm_solicitation_repo.NewPgCrmSolicitationRepository(conn, logger)
	httpClient := &server.NetHttpClient{}

	return usecases.NewCrmExportUseCase(httpClient, &adapters.HTTPDownloader{}, specRepo, companyRepo, solicitationRepo, mailer, logger)
}

func handleExportRequest(d amqp.Delivery, conn *pgxpool.Pool, client *messaging.RabbitClient) {
	ctx := getMessageContext(d)
	defer func(ctx context.Context) {
		tx := apm.TransactionFromContext(ctx)
		if tx != nil {
			tx.End()
		}
	}(ctx)

	logger := mainLogger.With(apmzap.TraceContext(ctx)...)
	logger.Info("Received message on exports queue", zap.Any("message", d))

	var req usecases.ExportRequest
	if err := json.Unmarshal(d.Body, &req); err != nil {
		logger.Error("Failed to unmarshal message", zap.Error(err))
		failOnError(d.Nack(false, false), "Failed to nack message")
	}

	sheetUc := getSheetUseCase(logger, conn)
	downloadUrl, err := sheetUc.Execute(req)
	publishResult(ctx, client, logger, req, downloadUrl, err)

	failOnError(d.Ack(false), "Failed to ack message")
}

func publishResult(ctx context.Context, c *messaging.RabbitClient, logger *zap.Logger, req usecases.ExportRequest, downloadUrl string, err error) {
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

	failOnError(c.Publish(ctx, "exports.results.excel", b), "Failed to publish response message")
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

func getSheetUseCase(logger *zap.Logger, conn *pgxpool.Pool) *usecases.SheetExportUseCase {
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

func getMessageContext(msg amqp.Delivery) context.Context {
	tp := getTraceParent(msg)
	if tp == "" {
		return context.Background()
	}

	t, err := apmhttp.ParseTraceparentHeader(tp)
	if err != nil {
		return context.Background()
	}

	tx := apm.DefaultTracer().StartTransactionOptions("Exporting Sheet", "message", apm.TransactionOptions{
		TraceContext: t,
	})

	return apm.ContextWithTransaction(context.Background(), tx)
}

func getTraceParent(msg amqp.Delivery) string {
	rawTp, ok := msg.Headers["traceparent"]
	if !ok {
		log.Println("traceparent not found in headers")
		return ""
	}
	tp, ok := rawTp.(string)
	if !ok {
		log.Println("traceparent is not a string")
		return ""
	}
	return tp
}

func getLogger() *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		os.Stdout,
		zap.DebugLevel,
	))
	return logger.With(zap.String("service.name", os.Getenv("ELASTIC_APM_SERVICE_NAME")), zap.String("service.environment", os.Getenv("ELASTIC_APM_ENVIRONMENT")))
}
