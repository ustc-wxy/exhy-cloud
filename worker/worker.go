package worker

import (
	"container/list"
	"exhy-cloud/task"
	"log"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     list.List
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	//toDo
}

func (w *Worker) RunTask() {
	//toDo
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config, "")
	result := d.Run()
	if result.Error != nil {
		log.Printf("[Worker] Error running task %v: %v", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}

	t.StartTime = time.Now().UTC()
	t.ContainerID = result.ContainerId
	t.State = task.Running
	w.Db[t.ID] = &t

	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config, t.ContainerID)
	result := d.Stop()
	if result.Error != nil {
		log.Printf("[Worker] Error stopping container %v: %v", t.ContainerID, result.Error)
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t

	return result
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.PushBack(t)
}
