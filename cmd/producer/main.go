package main

import (
	"context"
	"encoding/json"
	"export-service/internal/messaging"
	"export-service/internal/usecases"
	"log"

	"go.uber.org/zap"
)

func main() {

	con, err := messaging.Connect("guest", "guest", "localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer con.Close()

	logger := zap.NewExample()
	client := messaging.NewRabbitMQClient(con, logger)

	defer client.Close()
	req := usecases.ExportRequest{
		UserEmail:       "victor@driva.com.br",
		UserCompany:     "Driva",
		DataSource:      "linkedin",
		DataDownloadURL: "https://1.1.1.1",
	}
	reqBytes, _ := json.Marshal(req)
	failOnError(client.Publish(context.Background(), "exports.excel", reqBytes), "Failed to publish message")
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
