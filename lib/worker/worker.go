// Package worker contains the worker logic and some basic types
package worker

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Worker represents worker data for handling tasks
type Worker struct {
	rabbitmqURL string             // URL for RabbitMQ connection
	handler     TaskHandler        // Task handler function. Worker will call this function when it receives a task
	queueName   string             // Name of the queue to be subscribed to
	logger      *zap.SugaredLogger // Logger to write to
	broker      Broker             // Broker to receive and send messages to
}

// Task represents a task for processing
type Task struct {
	Payload *string // Data for processing
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
	var err error
	err = w.broker.Connect()

	w.fatalOnFail(err, "Failed to connect to RabbitMQ")

	receiving, err := w.broker.Subscribe(w.queueName)

	w.fatalOnFail(err, "Failed to subscribe to RabbitMQ queue")

	var forever chan struct{}

	w.logger.Infof("Worker starting...")
	go func() {
		for d := range receiving {
			w.logger.Info("Received message")
			body := string(d.Body)
			w.logger.Debug(body)
			err := w.handler(HandlerContext{Task: Task{&body}, broker: w.broker, Logger: w.logger})
			if err != nil {
				w.logger.Warnf("Failed to process the task: %s", err.Error())
				err = d.Reject(false) // TODO think about this behavior
				if err != nil {
					w.logger.Errorf("Failed to reject: %s", err.Error())
				}
				continue
			}
			err = d.Ack(false)
			if err != nil {
				w.logger.Errorf("Failed to send ACK: %s", err.Error())
				return
			}
		}
	}()

	return forever
}

// Stop the active connections
func (w *Worker) Stop() {
	err := w.broker.Close()
	if err != nil {
		return
	}
}

// New function creates new worker instance
func New(rabbitmqURL string, queueName string, handler TaskHandler) *Worker {
	broker := &RabbitMQBroker{rabbitmqURL: rabbitmqURL}
	return &Worker{
		rabbitmqURL: rabbitmqURL,
		handler:     handler,
		logger:      CreateDefaultLogger(zapcore.InfoLevel), // TODO read from config
		queueName:   queueName,
		broker:      broker,
	}
}

// CreateDefaultLogger creates default logger with a certain log level
func CreateDefaultLogger(level zapcore.Level) *zap.SugaredLogger {
	cfg := zap.Config{
		Encoding:         "console",
		Development:      true,
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:          "message",
			LevelKey:            "level",
			TimeKey:             "name",
			NameKey:             "ts",
			CallerKey:           "caller",
			FunctionKey:         "func",
			StacktraceKey:       "stacktrace",
			SkipLineEnding:      false,
			LineEnding:          "\n",
			EncodeLevel:         zapcore.CapitalColorLevelEncoder,
			EncodeTime:          zapcore.ISO8601TimeEncoder,
			EncodeDuration:      zapcore.MillisDurationEncoder,
			EncodeCaller:        zapcore.FullCallerEncoder,
			EncodeName:          zapcore.FullNameEncoder,
			NewReflectedEncoder: nil,
			ConsoleSeparator:    "\t",
		},
	}
	logger, _ := cfg.Build()
	return logger.Sugar()
}
