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

// Publisher is used to send messages to RabbitMQ.
type Publisher struct {
	addr    string
	queue   string
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
}

func NewPublisher(addr, queue string) *Publisher {
	return &Publisher{
		addr:  addr,
		queue: queue,
		done:  make(chan error),
	}
}

// Send puts data to queue.
func (p *Publisher) Send(data []byte) error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}
		<-time.After(d)
		if err := p.Connect(); err != nil {
			zap.L().Debug("could not reconnect", zap.Error(err))

			continue
		}
		err := p.channel.Publish(
			"",
			p.queue,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         data,
			})
		if err != nil {
			zap.L().Error("failed to send data", zap.Error(err))

			continue
		}
		zap.L().Debug("sent", zap.String("data", string(data)))

		return nil
	}
}

// Connects initialises Publisher.
func (p *Publisher) Connect() error {
	conn, err := amqp.Dial(p.addr)
	if err != nil {
		return err
	}
	p.conn = conn

	p.channel, err = p.conn.Channel()
	if err != nil {
		return err
	}

	go func() {
		zap.L().Debug("closing", zap.Error(<-p.conn.NotifyClose(make(chan *amqp.Error))))
		p.done <- errors.New("channel closed")
	}()

	_, err = p.channel.QueueDeclare(
		p.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

// Close stops Publisher.
func (p *Publisher) Close() error {
	//err := p.channel.Close()
	//if err != nil {
	//return err
	//}

	return p.conn.Close()
}
