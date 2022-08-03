package main

import "github.com/codex-team/hawk.workers.go/lib/worker"

func main() {
	rabbitmqUrl := "amqp://rabbitmq"
	workerInstance := worker.New(rabbitmqUrl, Handler)

	workerInstance.Run()
}
