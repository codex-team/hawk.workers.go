package worker

import (
	"context"

	"github.com/codex-team/hawk.workers.go/pkg/rmq"
	"github.com/streadway/amqp"
	"github.com/valyala/fastjson"
	"go.uber.org/zap"
)

// Handler is function that processes events and returns modified data.
type Handler func(*fastjson.Value) ([]byte, error)

// Worker reads data from queue, processes then and sends to another queue.
type Worker struct {
	pub     *rmq.Publisher
	con     *rmq.Consumer
	events  chan amqp.Delivery
	parser  fastjson.Parser
	handler Handler
}

func New(address, consumerQueue, publisherQueue string, h Handler) *Worker {
	events := make(chan amqp.Delivery)
	return &Worker{
		pub:     rmq.NewPublisher(address, publisherQueue),
		con:     rmq.NewConsumer(address, consumerQueue, events),
		events:  events,
		parser:  fastjson.Parser{},
		handler: h,
	}
}

// Run performs Worker logic.
func (w *Worker) Run(ctx context.Context) error {
	errs := make(chan error, 1)
	go func() {
		errs <- w.con.Receive(ctx)
	}()

	var data []byte
	for {
		select {
		case d := <-w.events:
			data = d.Body
			val, err := w.parser.ParseBytes(data)
			if err != nil {
				zap.L().Error("failed to parse data", zap.Error(err))
				continue
			}

			res, err := w.handler(val)
			if err != nil {
				zap.L().Error("failed to process data", zap.Error(err))
				continue
			}

			err = w.pub.Send(res)
			if err != nil {
				zap.L().Error("failed to send data", zap.Error(err))
				continue
			}

			err = d.Ack(false)
			if err != nil {
				zap.L().Error("failed to ack message", zap.Error(err))
			}
		case err := <-errs:
			if err != nil {
				zap.L().Error("worker disconnected with error", zap.Error(err))
			}
			return w.Stop()
		case <-ctx.Done():
			return w.Stop()
		}
	}
}

// Stop closes Publisher and Consumer.
func (w *Worker) Stop() error {
	err := w.pub.Close()
	if err != nil {
		zap.L().Error("publisher exited with error", zap.Error(err))
	}
	return w.con.Close()
}
