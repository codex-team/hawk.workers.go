package main

import "github.com/codex-team/hawk.workers.go/lib/worker"
import "encoding/json"
import "log"

func Handler(ctx worker.HandlerContext) error {
	var payload worker.Event
	err := json.Unmarshal([]byte(ctx.Task.Payload), &payload)
	if err != nil {
		log.Printf("Failing the task")
		return err
	}
	log.Printf("Resending the task")

	go ctx.SendTask(&ctx.Task)
	return nil
}
