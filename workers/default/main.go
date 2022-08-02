package worker_default

import "github.com/codex-team/hawk.workers.go/lib"

func main() {
	rabbitmqUrl := "amqp://rabbitmq"
	worker := lib.NewWorker(rabbitmqUrl, Handler)

	worker.Run()
}
