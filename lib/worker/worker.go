// Package worker contains the worker logic and some basic types
package worker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// Worker represents worker data for handling tasks
type Worker struct {
	rabbitmqURL string           // URL for RabbitMQ connection
	handler     TaskHandler      // Task handler function. Worker will call this function when it receives a task
	queueName   string           // Name of the queue to be subscribed to
	logger      *log.Logger      // Logger to write to
	connection  *amqp.Connection // Connection to the RabbitMQ
	channel     *amqp.Channel    // Channel to the RabbitMQ
}

// HandlerContext will be passed to the handler function on every call
type HandlerContext struct {
	Task    Task          // Task for processing
	channel *amqp.Channel // Channel to which it is connected to
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

// SendTask sends task to another queue, empty string is considered as no-op and nothing will be sent
func (ctx *HandlerContext) SendTask(task *Task, queueName string) error {
	if len(queueName) == 0 {
		return nil // considered as is not intended to be resent
	}
	err := ctx.channel.PublishWithContext(
		context.TODO(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(task.Payload),
		})

	if err != nil {
		log.Printf("Failed to send to another queue: %s", err.Error())
		return err
	}
	return nil
}

// Run function starts the worker
func (w *Worker) Run() <-chan struct{} {
	var err error
	w.connection, err = amqp.Dial(w.rabbitmqURL)
	w.fatalOnFail(err, "Failed to connect to RabbitMQ")

	w.channel, err = w.connection.Channel()
	w.fatalOnFail(err, "Failed to get channel to RabbitMQ")

	receiving, err := w.channel.Consume(
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
			err := w.handler(HandlerContext{Task: Task{string(d.Body)}, channel: w.channel})
			if err != nil {
				w.logger.Printf("Error on processing the task: %s", err.Error())
				err = d.Reject(false) // TODO think about this behavior
				if err != nil {
					w.logger.Printf("Failed to reject: %s", err.Error())
				}
				continue
			}
			err = d.Ack(false)
			if err != nil {
				w.logger.Printf("Failed to send ACK: %s", err.Error())
				return
			}
		}
	}()

	return forever
}

// Stop the active connections
func (w *Worker) Stop() {
	if w.channel != nil {
		err := w.channel.Close()
		if err != nil {
			w.logger.Printf("Failed to close the channel: %s", err.Error())
		}
	}
	if w.connection != nil {
		err := w.connection.Close()
		if err != nil {
			w.logger.Printf("Failed to close the connection: %s", err.Error())
		}
	}
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
