package messaging

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQExportRequestListener struct {
	c *amqp.Channel

	logger *zap.Logger
}

func NewRabbitMQExportRequestListener(c *amqp.Channel, logger *zap.Logger) *RabbitMQExportRequestListener {

	return &RabbitMQExportRequestListener{
		c: c,

		logger: logger.Named("RabbitMQExportRequestListener"),
	}
}

func (listener *RabbitMQExportRequestListener) StartConsuming(ctx context.Context) error {
	q, err := listener.c.QueueDeclare("export-requests", true, false, false, false, nil)

	if err != nil {
		listener.logger.Error("Failed to declare a queue", zap.Error(err))
		return err
	}

	msgs, err := listener.c.Consume(q.Name, "exports-consumer", true, false, false, false, nil)

	if err != nil {
		listener.logger.Error("Failed to consume messages", zap.Error(err))
		return err
	}

	listener.logger.Info("Waiting for messages")
	for {
		select {
		case <-ctx.Done():
			listener.logger.Info("Stopped consuming messages because of Done channel")
			return ctx.Err()
		case d := <-msgs:
			listener.logger.Info("Consumed message", zap.Any("message", string(d.Body)))
			// TODO: process message
		default:
			time.Sleep(time.Second)
		}

	}

}
