package main

import (
	"encoding/json"
	"github.com/codex-team/hawk.workers.go/lib/worker"
	"log"
)

const targetQueue string = "grouper"

func Handler(ctx worker.HandlerContext) error {
	var payload worker.Event
	err := json.Unmarshal([]byte(ctx.Task.Payload), &payload)
	if err != nil {
		log.Printf("Failing the task")
		return err
	}
	log.Printf("Resending the task")

	return ctx.SendTask(&ctx.Task, targetQueue)
}
