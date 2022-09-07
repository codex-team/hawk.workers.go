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

	// todo: read config
	// todo: setup connection to MongoDB database

	handler := CreateHandler(GrouperDeps{"eventsDb", "redis"}) // init db connection and pass it to handler

	workerInstance := worker.New(rabbitmqUrl, queue, handler)
	defer workerInstance.Stop() // TODO gracefully close connections on exit
	<-workerInstance.Run()
}
