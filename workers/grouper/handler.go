package main

import (
	"github.com/codex-team/hawk.workers.go/lib/worker"
)

// GrouperDeps is a struct with dependencies for grouper handler
type GrouperDeps struct {
	eventsDb string // todo: change to real type
	redis    string // todo: change to real type
}

// CreateHandler creates grouper handler with given deps
func CreateHandler(d GrouperDeps) func(ctx worker.HandlerContext) error {
	return func(ctx worker.HandlerContext) error {
		ctx.Logger.Debug(d.eventsDb)
		return nil
	}
}
