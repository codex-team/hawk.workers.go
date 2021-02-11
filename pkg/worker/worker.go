package worker

import (
	"context"
	"log"

	"github.com/codex-team/hawk.workers.go/pkg/rmq"
	"github.com/valyala/fastjson"
)

// Handler is function that processes events and returns modified data.
type Handler func(*fastjson.Value) ([]byte, error)

// Worker reads data from queue, processes then and sends to another queue.
type Worker struct {
	pub     *rmq.Publisher
	con     *rmq.Consumer
	events  chan []byte
	parser  fastjson.Parser
	handler Handler
}

func New(address, consumerQueue, publisherQueue string, h Handler) *Worker {
	events := make(chan []byte)
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

	for {
		select {
		case data := <-w.events:
			val, err := w.parser.ParseBytes(data)
			if err != nil {
				log.Printf("failed to parse data: %s", err.Error())
				continue
			}
			res, err := w.handler(val)
			if err != nil {
				err = w.pub.Send(res)
				if err != nil {
					log.Printf("failed to send data: %s", err.Error())
				}
			} else {
				log.Printf("failed to process data: %s", err.Error())
			}
		case err := <-errs:
			if err != nil {
				log.Printf("worker disconnected with error: %s", err.Error())
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
		log.Printf("publisher exited with error: %s", err.Error())
	}
	return w.con.Close()
}
