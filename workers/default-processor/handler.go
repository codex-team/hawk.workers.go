package main

import "github.com/codex-team/hawk.workers.go/lib/worker"

func Handler(ctx worker.HandlerContext) error {
	ctx.SendTask(&worker.Task{Payload: "Hello World!"})
	return nil
}
