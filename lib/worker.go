package lib

type Worker struct {
	rabbitmqUrl string
	handler     WorkerHandler
}

type WorkerContext struct {
	task    Task
	AddTask func(task *Task)
}

type Task struct {
	Payload string
}

func (w *Worker) Run() {

}

type WorkerHandler func(ctx WorkerContext) error

func NewWorker(rabbitmqUrl string, handler WorkerHandler) *Worker {
	return &Worker{rabbitmqUrl: rabbitmqUrl, handler: handler}
}
