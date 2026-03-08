package rabbitmq

import (
	"context"
	"discordcommandbot/internal/config"
	"discordcommandbot/pkg/errors"
	"discordcommandbot/pkg/logger"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	ctx    context.Context
	Config *config.RabbitMQConfig
	Conn   *amqp.Connection
	Ch     *amqp.Channel
}

func NewRabbitMQIntegration(ctx context.Context, config *config.RabbitMQConfig) (*RabbitMQ, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, errors.NewIntegration("failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, errors.NewIntegration("failed to open a channel", err)
	}

	rmq := &RabbitMQ{
		ctx:    ctx,
		Config: config,
		Conn:   conn,
		Ch:     ch,
	}

	if err := rmq.setup(); err != nil {
		rmq.Close()
		return nil, err
	}

	return rmq, nil
}

func (r *RabbitMQ) setup() error {
	err := r.Ch.ExchangeDeclare(
		r.Config.ExchangeName, // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return errors.NewIntegration("failed to declare exchange", err)
	}

	_, err = r.Ch.QueueDeclare(
		r.Config.QueueName, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return errors.NewIntegration("failed to declare queue", err)
	}

	err = r.Ch.QueueBind(
		r.Config.QueueName,
		r.Config.RoutingKey,
		r.Config.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return errors.NewIntegration("failed to bind queue", err)
	}

	logger.Info("✅ RabbitMQ successfully configured")
	return nil
}

func (r *RabbitMQ) Publisher(body []byte) error {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	err := r.Ch.PublishWithContext(
		ctx,
		r.Config.ExchangeName, // exchange
		r.Config.RoutingKey,   // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return errors.NewIntegration("failed to publish message", err)
	}

	logger.Debug("📨 Message published: %s", string(body))
	return nil
}

func (r *RabbitMQ) Consumer() (<-chan amqp.Delivery, error) {

	msgs, err := r.Ch.Consume(
		r.Config.QueueName, // queue
		"go-consumer",      // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return nil, errors.NewIntegration("failed to register consumer", err)
	}
	return msgs, nil
}

func (r *RabbitMQ) Close() error {
	if r.Ch != nil {
		if err := r.Ch.Close(); err != nil {
			return err
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
