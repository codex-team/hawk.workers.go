// Package worker contains the worker logic and some basic types
package worker

import "fmt"

// Worker represents worker data for handling tasks
type Worker struct {
	rabbitmqURL string      // URL for RabbitMQ connection
	handler     TaskHandler // Task handler function. Worker will call this function when it receives a task
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

// Run function starts the worker
func (w *Worker) Run() {
	// @todo rabbitmq connection login here
	fmt.Println("Worker started...")
}

// New function creates new worker instance
func New(rabbitmqURL string, handler TaskHandler) *Worker {
	return &Worker{rabbitmqURL: rabbitmqURL, handler: handler}
}
