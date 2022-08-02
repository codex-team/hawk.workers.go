package worker_default

import "github.com/codex-team/hawk.workers.go/lib"

func Handler(ctx lib.WorkerContext) error {
	ctx.AddTask(&lib.Task{Payload: "Hello World!"})
	return nil
}
