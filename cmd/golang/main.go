package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/codex-team/hawk.workers.go/internal/workers/golang"
	"github.com/codex-team/hawk.workers.go/pkg/worker"
)

const write_queue = "grouper"

var (
	rabbit_addr string
	read_queue  string
)

func init() {
	flag.StringVar(&rabbit_addr, "r", "amqp://127.0.0.1:5672", "RabbitMQ address")
	flag.StringVar(&read_queue, "q", "errors/golang", "Queue to read from")
}

func main() {
	flag.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker := worker.New(rabbit_addr, read_queue, write_queue, golang.Handler)
	go func() {
		err := worker.Run(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	signal.Stop(done)
}
