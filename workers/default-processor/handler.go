package main

import "github.com/codex-team/hawk.workers.go/lib/worker"
import "encoding/json"

func Handler(ctx worker.HandlerContext) error {
	var payload worker.Event
	err := json.Unmarshal([]byte(ctx.Task.Payload), &payload)
	if err != nil {
		return err
	}

	go ctx.SendTask(&ctx.Task)
	return nil
}
