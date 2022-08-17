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

// TestHandler tests Handler for correct error handling and resending the task
func TestHandler(t *testing.T) {
	logger, logs := setupLogsCapture()
	contextStamp, currentTask := 0, ""
	validJson := `{"projectId":"sample","payload":{"title":"payload title","type":"sample error","backtrace":[{"file":"main.js","line":42,"column":42,"function":null,"sourceCode":[{"line":42,"content":"some content"}]}],"context":{},"catcherVersion":"0.1.0","timestamp":1659587713},"catcherType":"errors/nodejs"}`
	sendTask := func(ctx *worker.HandlerContext, task *worker.Task, queueName string) error {
		contextStamp = 42
		currentTask = *task.Payload
		if queueName == "grouper" {
			return nil
		}
		return errors.New("queueName is not grouper")
	}
	err := Handler(worker.HandlerContext{Task: worker.Task{Payload: &validJson}, Logger: logger.Sugar(), SendTask: sendTask})
	if err != nil ||
		logs.Len() != 1 ||
		logs.All()[0].Level != zap.DebugLevel ||
		logs.All()[0].Message != "Resending the task" ||
		contextStamp != 42 ||
		currentTask != validJson {
		t.Error("Failed valid JSON test")
	}

	contextStamp, currentTask = 0, ""
	logs.TakeAll()
	err = Handler(worker.HandlerContext{Task: worker.Task{Payload: nil}, Logger: logger.Sugar(), SendTask: sendTask})
	if err.Error() != "Task.Payload is nil" ||
		logs.Len() != 1 ||
		logs.All()[0].Level != zap.ErrorLevel ||
		logs.All()[0].Message != "Error in the context: Task.Payload is nil" ||
		contextStamp != 0 {
		t.Error("Failed null payload test")
	}

	contextStamp, currentTask = 0, ""
	logs.TakeAll()
	invalidJson := `{"projectId": "esf", "payload":{"type": 777}}`
	err = Handler(worker.HandlerContext{Task: worker.Task{Payload: &invalidJson}, Logger: logger.Sugar(), SendTask: sendTask})
	if err.Error() != "json: cannot unmarshal number into Go struct field Payload.payload.type of type string" ||
		logs.Len() != 1 ||
		logs.All()[0].Level != zap.DebugLevel ||
		logs.All()[0].Message != "Failing the task" ||
		contextStamp != 0 {
		t.Error("Failed invalid JSON structure payload test")
	}

	contextStamp, currentTask = 0, ""
	logs.TakeAll()
	noprojectJson := `{"some": "esf", "payload":{"type": ""}}`
	err = Handler(worker.HandlerContext{Task: worker.Task{Payload: &noprojectJson}, Logger: logger.Sugar(), SendTask: sendTask})
	if err.Error() != "no projectId or no payload.type fields" ||
		logs.Len() != 1 ||
		logs.All()[0].Level != zap.WarnLevel ||
		logs.All()[0].Message != "Fail the task: no projectId or no payload.type fields" ||
		contextStamp != 0 {
		t.Error("Failed no-project-JSON structure payload test")
	}

}
