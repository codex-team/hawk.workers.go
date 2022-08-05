package main

import "github.com/codex-team/hawk.workers.go/lib/worker"

func main() {
	rabbitmqUrl := "amqp://127.0.0.1:5672"
	queue := "errors/default"
	workerInstance := worker.New(rabbitmqUrl, queue, Handler, nil)

	<-workerInstance.Run()
}
