package worker

import (
	"container/list"
	"exhy-cloud/task"

	"github.com/google/uuid"
)

type Worker struct {
	Queue     list.List
	Db        map[uuid.UUID]task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	//toDo
}

func (w *Worker) RunTask() {
	//toDo
}

func (w *Worker) StartTask() {
	//toDo
}

func (w *Worker) StopTask() {
	//toDo
}
