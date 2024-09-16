package main

import (
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

	// Create queues
	exports := "exports.excel"
	failOnError(client.CreateQueue(exports, true), "Failed to create exports queue")

	dlq := "exports.dlq"
	failOnError(client.CreateQueue(dlq, false), "Failed to create dlq queue")

	// Bind DLQ
	failOnError(client.BindDLQ(dlq), "Failed to bind DLQ")

	exportsBus, err := client.Consume(exports)
	failOnError(err, "Failed to consume bus")

	dlqBus, err := client.Consume(dlq)
	failOnError(err, "Failed to consume dlq bus")

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
		}
	}()

	go func() {
		for d := range dlqBus {
			logger.Info("Received message on dlq", zap.Any("message", d))
			if d.Redelivered {
				// Ignores redelivered dlq messages
				logger.Debug("Redelivered message", zap.Any("message", d))
				continue
			}

			// Unmarshal Body into DTO
			var req usecases.ExportRequest
			if err = json.Unmarshal(d.Body, &req); err != nil {
				logger.Error("Failed to unmarshal message", zap.Error(err))
				failOnError(d.Nack(false, true), "Failed to nack message")
			}
		}
	}()

	logger.Info("Consuming messages, press CTRL+C to stop")
	// Blocks forever
	<-blocking
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
