package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/codex-team/hawk.workers.go/lib/worker"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const targetQueue = "grouper"

func Handler(ctx worker.HandlerContext) error {
	if ctx.Channel == nil {
		return errors.New("ctx.Channel is nil")
	}
	var payload worker.Event
	err := json.Unmarshal([]byte(ctx.Task.Payload), &payload)
	if err != nil {
		log.Printf("Failing the task")
		return err
	}
	log.Printf("Resending the task")

	err = ctx.Channel.PublishWithContext(
		context.TODO(),
		"",
		targetQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(ctx.Task.Payload),
		})

	if err != nil {
		log.Printf("Failed to send to another queue: %s", err.Error())
		return err
	}
	ctx.SendTask(&ctx.Task)
	return nil
}
