package messaging

import (
	"context"
	"fmt"
	"time"

	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitClient struct {
	ch     *amqp.Channel
	logger *zap.Logger
}

func Connect(username, password, host string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", username, password, host))
}
func NewRabbitMQClient(conn *amqp.Connection, logger *zap.Logger) *RabbitClient {
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("failed to open a channel", zap.Error(err))
	}
	return &RabbitClient{
		ch:     ch,
		logger: logger.Named("RabbitMQConnection"),
	}
}

func (c *RabbitClient) CreateQueue(name string, dlxName *string, routingKey *string) error {
	if name == "" {
		return fmt.Errorf("queue name cannot be empty")
	}
	args := amqp.Table{}
	if dlxName != nil && routingKey != nil {
		args = amqp.Table{"x-dead-letter-exchange": *dlxName,
			"x-dead-letter-routing-key": *routingKey,
		}
	}
	_, err := c.ch.QueueDeclare(
		name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)
	return err
}

func (c *RabbitClient) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(
		queue,
		"",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		amqp.Table{
			"x-priority": int(time.Now().Unix()),
		}, // args
	)
}

func (c *RabbitClient) Publish(ctx context.Context, queue string, body []byte) error {
	return c.ch.Publish(
		"",
		queue,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers: amqp.Table{
				"traceparent": apmhttp.FormatTraceparentHeader(apm.TransactionFromContext(ctx).TraceContext()),
			},
		},
	)
}

func (c *RabbitClient) BindDLQ(queue string) error {
	return c.ch.QueueBind(
		queue,
		"",
		"exports-dlx",
		false,
		nil,
	)
}

func (c *RabbitClient) Close() error {
	return c.ch.Close()
}
