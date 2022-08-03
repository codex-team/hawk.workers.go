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
	queue_name  string
	logger      *log.Logger
}

// HandlerContext will be passed to the handler function on every call
type HandlerContext struct {
	Task     Task             // Task for processing
	SendTask func(task *Task) // Function for sending task to another worker
}

// Task represents a task for processing
type Task struct {
	Payload string // Data for processing
}

// TaskHandler represents Worker handler for processing tasks
type TaskHandler func(ctx HandlerContext) error

func (w *Worker) fatalOnFail(err error, msg string) {
	if err != nil {
		w.logger.Fatalf("%s: %s\n", msg, err.Error())
		return
	}
}

// Run function starts the worker
func (w *Worker) Run() <-chan struct{} {
	conn, err := amqp.Dial(w.rabbitmqURL)
	w.fatalOnFail(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	w.fatalOnFail(err, "Failed to get channel to RabbitMQ")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		w.queue_name, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	w.fatalOnFail(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	var forever chan struct{}

	w.logger.Println("Worker starting...")
	go func() {
		for d := range msgs {
			w.logger.Printf("Received message: %s", d.Body)
			w.handler(HandlerContext{Task: Task{string(d.Body)}, SendTask: func(task *Task) {}}) // TODO make proper SendTask logic
		}
	}()

	return forever
}

// New function creates new worker instance
func New(rabbitmqURL string, queue_name string, handler TaskHandler, logger *log.Logger) *Worker {
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
		queue_name: queue_name,
	}
}
