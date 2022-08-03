package main

import "github.com/codex-team/hawk.workers.go/lib/worker"

func main() {
	rabbitmqUrl := "amqp://rabbitmq"
	queue := "hello"
	workerInstance := worker.New(rabbitmqUrl, queue, Handler, nil)

	<-workerInstance.Run()
}
