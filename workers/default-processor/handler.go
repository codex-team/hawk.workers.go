package main

import (
	"encoding/json"
	"github.com/codex-team/hawk.workers.go/lib/worker"
)

const targetQueue string = "grouper"

func Handler(ctx worker.HandlerContext) error {
	var payload worker.Event
	err := json.Unmarshal([]byte(ctx.Task.Payload), &payload)
	if err != nil {
		ctx.Logger.Debug("Failing the task")
		return err
	}
	ctx.Logger.Debug("Resending the task")
	return ctx.SendTask(&ctx.Task, targetQueue)
}
