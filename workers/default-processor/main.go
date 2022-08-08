package main

import (
	"github.com/codex-team/hawk.workers.go/lib/worker"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	rabbitmqUrl := getEnv("REGISTRY_URL", "amqp://127.0.0.1:5672")
	queue := "errors/default"
	workerInstance := worker.New(rabbitmqUrl, queue, Handler, nil)

	<-workerInstance.Run()
}
