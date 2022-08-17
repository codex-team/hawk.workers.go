package main

import (
	"encoding/json"
	"errors"
	"github.com/codex-team/hawk.workers.go/lib/worker"
)

const targetQueue string = "grouper"

func Handler(ctx worker.HandlerContext) error {
	if ctx.Task.Payload == nil {
		ctx.Logger.Error("Error in the context: Task.Payload is nil")
		return errors.New("Task.Payload is nil")
	}
	var payload worker.Event
	err := json.Unmarshal([]byte(*ctx.Task.Payload), &payload)
	if err != nil {
		ctx.Logger.Debug("Failing the task")
		return err
	}
	if len(payload.ProjectId) == 0 || len(payload.Payload.Type) == 0 {
		ctx.Logger.Warn("Fail the task: no projectId or no payload.type fields")
		return errors.New("no projectId or no payload.type fields")
	}
	ctx.Logger.Debug("Resending the task")
	return ctx.SendTask(&ctx, &ctx.Task, targetQueue)
}
