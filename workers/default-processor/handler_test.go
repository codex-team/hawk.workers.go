package main_test

import (
	"github.com/codex-team/hawk.workers.go/lib/worker"
	"github.com/codex-team/hawk.workers.go/workers/default-processor"
	"github.com/stretchr/testify/assert"
	"mocks"
	"testing"
)

func TestHandler(t *testing.T) {
	payload := `
{
  "projectId": "62160d67a01657bf20b4627d",
  "payload": {
    "title": "TypeError: Cannot read properties of undefined (reading 'MemberId')",
    "type": "TypeError",
    "backtrace": [],
    "context": {},
    "catcherVersion": "3.1.2",
    "timestamp": 1659587713
  },
  "catcherType": "errors/nodejs"
}
`
	task := worker.Task{
		Payload: &payload,
	}

	expectedPayload := `
{
  "projectId": "62160d67a01657bf20b4627d",
  "payload": {
    "title": "TypeError: Cannot read properties of undefined (reading 'MemberId')",
    "type": "TypeError",
    "backtrace": [],
    "context": {},
    "catcherVersion": "3.1.2",
    "timestamp": 1659587713
  },
  "catcherType": "errors/nodejs"
}
`
	brokerMock := mocks.NewBroker(t)

	brokerMock.On("Publish", "grouper", []byte(expectedPayload)).Return(nil)

	ctx := worker.CreateHandlerContext(task, worker.CreateDefaultLogger(), brokerMock)

	err := main.Handler(ctx)

	assert.Nil(t, err)
}
