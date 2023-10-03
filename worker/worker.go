package worker

import (
	"container/list"
	"errors"
	"exhy-cloud/task"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     *list.List
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	//toDo
}

func (w *Worker) RunTask() task.DockerResult {
	t := w.Queue.Front().Value
	if t == nil {
		log.Println("[Worker] No tasks in the queue")
		return task.DockerResult{Error: nil}
	}
	w.Queue.Remove(w.Queue.Front())

	taskQueued := t.(task.Task)

	taskPersisted := w.Db[taskQueued.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = &taskQueued
	}

	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result.Error = errors.New("Should not reach here")
		}
	} else {
		err := fmt.Errorf("Invalid transtion from %v to %v",
			taskPersisted.State, taskQueued.State)
		result.Error = err
	}
	return result
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
