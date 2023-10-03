package worker

import (
	"container/list"
	"exhy-cloud/task"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestWorkerBasic(t *testing.T) {
	db := make(map[uuid.UUID]*task.Task)
	queue := list.New()
	worker := Worker{
		Queue: queue,
		Db:    db,
	}

	task_ := task.Task{
		ID:    uuid.New(),
		Name:  "test-task-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	fmt.Println("[Test] starting task ...")
	worker.AddTask(task_)

	result := worker.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
	task_.ContainerID = result.ContainerId

	fmt.Printf("[Test] task %s is running in container %s\n", task_.ID, task_.ContainerID)
	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 20)

	fmt.Println("[Test] stopping task ...")
	task_.State = task.Completed

	worker.AddTask(task_)

	result = worker.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

}
