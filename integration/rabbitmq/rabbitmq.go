package rabbitmq

import (
	"enque-learning/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Config config.Config
}

func NewRabbitMQIntegration(config config.Config) *RabbitMQ {
	return &RabbitMQ{
		Config: config,
	}
}

func (r *RabbitMQ) openChannel() (*amqp.Channel, error) {

	conn, err := amqp.Dial(r.Config.RabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (r *RabbitMQ) Publisher(body []byte) error {
	channel, err := r.openChannel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.Publish(
		r.Config.ExchangeName, // exchange
		r.Config.QueueName,    // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Consumer() (<-chan amqp.Delivery, error) {
	channel, err := r.openChannel()
	if err != nil {
		return nil, err
	}

	return channel.Consume(
		r.Config.QueueName, // queue
		"go-consumer",      // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
}
