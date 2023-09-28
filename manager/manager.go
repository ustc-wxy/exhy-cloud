package manager

import (
	"container/list"
	"exhy-cloud/task"

	"github.com/google/uuid"
)

type Manager struct {
	Pending       list.List
	TaskDb        map[string][]task.Task
	EventDb       map[string][]task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {
	//toDo
}

func (m *Manager) UpdateTasks() {
	//toDo
}

func (m *Manager) SendWork() {
	//toDo
}
