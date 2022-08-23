package broker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DeliveryChannel <-chan amqp.Delivery

type Publisher interface {
	Publish(queueName string, payload []byte) error
}

type Broker interface {
	Publisher
	Connect() error
	Subscribe(queueName string) (DeliveryChannel, error)
	Close() error
}

type RabbitMQ struct {
	RabbitmqURL string           // URL for RabbitMQ connection
	connection  *amqp.Connection // Connection to the RabbitMQ
	channel     *amqp.Channel    // Channel to the RabbitMQ
}

func (b *RabbitMQ) Connect() error {
	var err error
	b.connection, err = amqp.Dial(b.RabbitmqURL)
	if err != nil {
		return err
	}
	b.channel, err = b.connection.Channel()
	if err != nil {
		return err
	}
	return nil
}

func (b *RabbitMQ) Subscribe(queueName string) (DeliveryChannel, error) {
	receiving, err := b.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	return receiving, err
}

func (b *RabbitMQ) Publish(queueName string, payload []byte) error {
	err := b.channel.PublishWithContext(
		context.TODO(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		})
	if err != nil {
		return err
	}
	return nil
}

func (b *RabbitMQ) Close() error {
	if b.channel != nil {
		err := b.channel.Close()
		if err != nil {
			return err
		}
	}
	if b.connection != nil {
		err := b.connection.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
