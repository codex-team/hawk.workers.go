package main

import (
	"errors"
	"github.com/codex-team/hawk.workers.go/lib/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"testing"
)

// setupLogsCapture creates logger, suitable for testing
func setupLogsCapture() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.DebugLevel)
	return zap.New(core), logs
}

type MockedHandlerContext struct {
	task   worker.Task        // Task for processing
	logger *zap.SugaredLogger // Logger to write to
	stamp  uint
}

func (ctx *MockedHandlerContext) SendTask(task *worker.Task, queueName string) error {
	ctx.stamp = 42
	ctx.task.Payload = task.Payload
	if queueName == "grouper" {
		return nil
	}
	return errors.New("queueName is not grouper")
}
func (ctx *MockedHandlerContext) Task() *worker.Task {
	return &ctx.task
}

func (ctx *MockedHandlerContext) Logger() *zap.SugaredLogger {
	return ctx.logger
}

// TestHandler tests Handler for correct error handling and resending the task
func TestHandler(t *testing.T) {
	logger, logs := setupLogsCapture()
	validJson := `{"projectId":"sample","payload":{"title":"payload title","type":"sample error","backtrace":[{"file":"main.js","line":42,"column":42,"function":null,"sourceCode":[{"line":42,"content":"some content"}]}],"context":{},"catcherVersion":"0.1.0","timestamp":1659587713},"catcherType":"errors/nodejs"}`
	ctx := MockedHandlerContext{task: worker.Task{Payload: &validJson}, logger: logger.Sugar()}
	err := Handler(&ctx)
	if err != nil ||
		logs.Len() != 1 ||
		logs.All()[0].Level != zap.DebugLevel ||
		logs.All()[0].Message != "Resending the task" ||
		ctx.stamp != 42 ||
		*ctx.task.Payload != validJson {
		t.Error("Failed valid JSON test")
	}

	ctx = MockedHandlerContext{task: worker.Task{Payload: nil}, logger: logger.Sugar()}
	logs.TakeAll()
	err = Handler(&ctx)
	if err.Error() != "Task.Payload is nil" ||
		logs.Len() != 1 ||
		logs.All()[0].Level != zap.ErrorLevel ||
		logs.All()[0].Message != "Error in the context: Task.Payload is nil" ||
		ctx.stamp != 0 {
		t.Error("Failed null payload test")
	}

	logs.TakeAll()
	invalidJson := `{"projectId": "esf", "payload":{"type": 777}}`
	ctx = MockedHandlerContext{task: worker.Task{Payload: &invalidJson}, logger: logger.Sugar()}
	err = Handler(&ctx)
	if err.Error() != "json: cannot unmarshal number into Go struct field Payload.payload.type of type string" ||
		logs.Len() != 1 ||
		logs.All()[0].Level != zap.DebugLevel ||
		logs.All()[0].Message != "Failing the task" ||
		ctx.stamp != 0 {
		t.Error("Failed invalid JSON structure payload test")
	}

	logs.TakeAll()
	noprojectJson := `{"some": "esf", "payload":{"type": ""}}`
	ctx = MockedHandlerContext{task: worker.Task{Payload: &noprojectJson}, logger: logger.Sugar()}
	err = Handler(&ctx)
	if err.Error() != "no projectId or no payload.type fields" ||
		logs.Len() != 1 ||
		logs.All()[0].Level != zap.WarnLevel ||
		logs.All()[0].Message != "Fail the task: no projectId or no payload.type fields" ||
		ctx.stamp != 0 {
		t.Error("Failed no-project-JSON structure payload test")
	}

}
