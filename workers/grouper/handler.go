package main

import (
	"github.com/codex-team/hawk.workers.go/lib/worker"
)

type GrouperDeps struct {
	eventsDb string
	redis    string
}

func CreateHandler(d GrouperDeps) func(ctx worker.HandlerContext) error {
	return func(ctx worker.HandlerContext) error {
		ctx.Logger.Debug(d.eventsDb)
		return nil
	}
}
