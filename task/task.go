package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)


type Task struct {
	ID            uuid.UUID
	ContainerID   string
	Name          string
	State         State
	Image         string
	Memory        int
	Disk          int
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

type TaskEvent struct {
	ID    uuid.UUID
	State State
	Task  Task
}

func NewConfig(t *Task) Config {
	return Config{
		Name:   t.Name,
		Image:  t.Image,
		Memory: int64(t.Memory),
		Disk:   int64(t.Disk),
	}
}
