package broker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

// DeliveryChannel represents a channel for receiving messages from a queue
type DeliveryChannel <-chan amqp.Delivery

// Publisher represents a publisher for publishing messages to a queue
type Publisher interface {
	Publish(queueName string, payload []byte) error // Publish sends a message to a queue
}

// Broker represents a broker for publishing and receiving messages from a queues
type Broker interface {
	Publisher

	Connect() error                                      // Connect establishes a connection to the RabbitMQ
	Subscribe(queueName string) (DeliveryChannel, error) // Subscribe subscribes to a queue and returns a channel for receiving messages
	Close() error                                        // Close terminates the connection to the RabbitMQ
}

// RabbitMQ is a broker implementation for RabbitMQ
type RabbitMQ struct {
	RabbitmqURL string           // URL for RabbitMQ connection
	connection  *amqp.Connection // Connection to the RabbitMQ
	channel     *amqp.Channel    // Channel to the RabbitMQ
}

// Connect establishes a connection to the RabbitMQ
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

// Subscribe subscribes to a queue and returns a channel for receiving messages
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

// Publish sends a message to a queue
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

// Close terminates the connection to the RabbitMQ
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
