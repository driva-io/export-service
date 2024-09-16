package messaging_test

import (
	"context"
	"testing"
	"time"

	messaging "export-service/internal/messaging/export_request_listener"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestRabbitMQExportRequestListener(t *testing.T) {
	t.Parallel()

	ch, close := getRabbitChannel(t)
	defer close()
	defer ch.Close()

	logger, _ := zap.NewProduction()
	listener := messaging.NewRabbitMQExportRequestListener(ch, logger)

	c, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err := listener.StartConsuming(c)

	require.NoError(t, err)

}

func getRabbitChannel(t *testing.T) (*amqp.Channel, func()) {
	t.Helper()

	ctx := context.Background()

	logger, _ := zap.NewProduction()

	rabbitmqContainer, err := rabbitmq.Run(ctx,
		"rabbitmq:3.12.12-management-alpine",
	)
	close := func() {
		if err := rabbitmqContainer.Terminate(ctx); err != nil {
			logger.Error("failed to terminate container", zap.Error(err))
		}
	}
	if err != nil {
		logger.Fatal("failed to start container", zap.Error(err))
	}
	url, _ := rabbitmqContainer.AmqpURL(ctx)

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Fatal("failed to connect to rabbitmq", zap.Error(err), zap.String("url", url))
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("failed to open a channel", zap.Error(err))
	}

	return ch, close
}
