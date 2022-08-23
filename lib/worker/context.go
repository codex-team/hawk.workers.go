package worker

import (
	"github.com/codex-team/hawk.workers.go/lib/broker"
	"go.uber.org/zap"
)

// HandlerContext will be passed to the handler function on every call
type HandlerContext struct {
	Task   Task               // Task for processing
	Logger *zap.SugaredLogger // Logger to write to
	broker broker.Publisher   // Channel to which it is connected to
}

// SendTask sends task to another queue, empty string is considered as no-op and nothing will be sent
func (ctx *HandlerContext) SendTask(task *Task, queueName string) error {
	if len(queueName) == 0 {
		return nil // considered as is not intended to be resent
	}
	err := ctx.broker.Publish(queueName, []byte(*task.Payload))

	if err != nil {
		ctx.Logger.Errorf("Failed to send to queue `%s`: %s", queueName, err.Error())
		return err
	}
	return nil
}
