// Package worker contains the worker logic and some basic types
package worker

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// Worker represents worker data for handling tasks
type Worker struct {
	rabbitmqURL string      // URL for RabbitMQ connection
	handler     TaskHandler // Task handler function. Worker will call this function when it receives a task
	queueName   string      // Name of the queue to be subscribed to
	logger      *log.Logger // Logger to write to
}

// HandlerContext will be passed to the handler function on every call
type HandlerContext struct {
	Task     Task             // Task for processing
	SendTask func(task *Task) // Function for sending task to another worker
	Channel  *amqp.Channel    // Channel, from which the Task came from
}

// Task represents a task for processing
type Task struct {
	Payload string // Data for processing
}

// TaskHandler represents Worker handler for processing tasks
type TaskHandler func(ctx HandlerContext) error

// Helper function for handling fatal errors
func (w *Worker) fatalOnFail(err error, msg string) {
	if err != nil {
		w.logger.Fatalf("%s: %s\n", msg, err.Error())
	}
}

// Run function starts the worker
func (w *Worker) Run() <-chan struct{} {
	conn, err := amqp.Dial(w.rabbitmqURL)
	w.fatalOnFail(err, "Failed to connect to RabbitMQ")

	//defer conn.Close()
	ch, err := conn.Channel()
	w.fatalOnFail(err, "Failed to get channel to RabbitMQ")
	//defer ch.Close()

	receiving, err := ch.Consume(
		w.queueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	var forever chan struct{}

	w.logger.Println("Worker starting...")
	go func() {
		for d := range receiving {
			w.logger.Printf("Received message: %s", d.Body) // TODO move under debug level
			err := w.handler(HandlerContext{Task: Task{string(d.Body)}, Channel: ch, SendTask: func(task *Task) {
				err := d.Ack(false)
				if err != nil {
					w.logger.Printf("Failed to send ACK: %s", err.Error())
					return
				}
			}})

			if err != nil {
				w.logger.Printf("Error on processing the task: %s", err.Error())
				err = d.Reject(false) // TODO think about this behavior
				if err != nil {
					w.logger.Printf("Failed to reject: %s", err.Error())
				}
				continue
			}
		}
	}()

	return forever
}

// New function creates new worker instance
func New(rabbitmqURL string, queueName string, handler TaskHandler, logger *log.Logger) *Worker {
	return &Worker{
		rabbitmqURL: rabbitmqURL,
		handler:     handler,
		logger: func() *log.Logger {
			if logger != nil {
				return logger
			} else {
				return log.Default()
			}
		}(),
		queueName: queueName,
	}
}
