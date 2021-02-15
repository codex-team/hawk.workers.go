package rmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// Consumer is used to read messages from RabbitMQ.
type Consumer struct {
	addr    string
	queue   string
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
	errs    chan amqp.Delivery
}

func NewConsumer(addr, queue string, errs chan amqp.Delivery) *Consumer {
	return &Consumer{
		addr:  addr,
		queue: queue,
		done:  make(chan error),
		errs:  errs,
	}
}

// Connect initialises Consumer.
func (c *Consumer) Connect() error {
	conn, err := amqp.Dial(c.addr)
	if err != nil {
		return err
	}
	c.conn = conn

	c.channel, err = c.conn.Channel()
	if err != nil {
		return err
	}

	go func() {
		zap.L().Debug("closing", zap.Error(<-c.conn.NotifyClose(make(chan *amqp.Error))))
		c.done <- errors.New("channel closed")
	}()

	_, err = c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

// Reconnect performs reconnecting to RabbitMQ using exponential backoff.
func (c *Consumer) Reconnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = 3 * time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 30 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		select {
		case <-ctx.Done():
			return nil, nil
		case <-time.After(d):
			if err := c.Connect(); err != nil {
				zap.L().Debug("could not connect in reconnect call", zap.Error(err))

				continue
			}
			msgs, err := c.channel.Consume(
				c.queue,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				zap.L().Debug("could not connect", zap.Error(err))

				continue
			}

			return msgs, nil
		}
	}
}

// Receive reads data from queue.
func (c *Consumer) Receive(ctx context.Context) error {
	msgs, err := c.Reconnect(ctx)
	if err != nil {
		return err
	}

	for {
		go func(msgs <-chan amqp.Delivery) {
			for d := range msgs {
				c.errs <- d
			}
		}(msgs)

		if <-c.done != nil {
			msgs, err = c.Reconnect(ctx)
			if err != nil {
				return err
			}
		}
	}
}

// Close stops connection.
func (c *Consumer) Close() error {
	err := c.channel.Close()
	if err != nil {
		return err
	}

	return c.conn.Close()
}
