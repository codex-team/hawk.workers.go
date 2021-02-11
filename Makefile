CMD := go-worker

build:
	go build -o $(CMD) ./cmd/golang/main.go

run:
	make build
	./$(CMD)

lint:
	golangci-lint run ./...

.PHONY: build run lint test
