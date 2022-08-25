package main_test

import (
	"github.com/codex-team/hawk.workers.go/lib/worker"
	"github.com/codex-team/hawk.workers.go/workers/default-processor"
	"github.com/stretchr/testify/assert"
	"mocks"
	"testing"
)

var logger = worker.CreateTestLogger()

func TestDefaultProcessor(t *testing.T) {
	t.Run("Should send correct task to grouper worker", func(t *testing.T) {
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

		brokerMock := mocks.NewBroker(t)
		brokerMock.On("Publish", "grouper", []byte(payload)).Return(nil)

		ctx := worker.CreateHandlerContext(task, logger, brokerMock)

		err := main.Handler(ctx)

		assert.Nil(t, err)
	})

	failedPayloads := map[string]string{
		"Should fail if there is no project id in the task": `
{
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
`,
		"Should fail if there is no payload in the task": `
{
  "projectId": "62160d67a01657bf20b4627d",
  "catcherType": "errors/nodejs"
}
`,
		"Should fail if there is no catherType in the task": `
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
}
`,
	}

	for name, payload := range failedPayloads {
		t.Run(name, func(t *testing.T) {
			task := worker.Task{
				Payload: &payload,
			}

			brokerMock := mocks.NewBroker(t)

			ctx := worker.CreateHandlerContext(task, logger, brokerMock)

			err := main.Handler(ctx)

			assert.Error(t, err)
		})
	}
}
