package worker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type DeliveryChannel <-chan amqp.Delivery

type BrokerPublisher interface {
	Publish(queueName string, payload []byte) error
}

type Broker interface {
	BrokerPublisher
	Connect() error
	Subscribe(queueName string) (DeliveryChannel, error)
	Close() error
}

type RabbitMQBroker struct {
	rabbitmqURL string           // URL for RabbitMQ connection
	connection  *amqp.Connection // Connection to the RabbitMQ
	channel     *amqp.Channel    // Channel to the RabbitMQ
}

func (b *RabbitMQBroker) Connect() error {
	var err error
	b.connection, err = amqp.Dial(b.rabbitmqURL)
	if err != nil {
		return err
	}
	b.channel, err = b.connection.Channel()
	if err != nil {
		return err
	}
	return nil
}

func (b *RabbitMQBroker) Subscribe(queueName string) (DeliveryChannel, error) {
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

func (b *RabbitMQBroker) Publish(queueName string, payload []byte) error {
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

func (b *RabbitMQBroker) Close() error {
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
