package task

import "github.com/docker/docker/client"

type Config struct {
	Name         string
	AttachStdin  bool
	AttachStdout bool
	AttachStderr bool
	Cmd          []string
	Image        string
	Memory       int64
	Disk         int64
	Env          []string
	ResultPolicy string
}

type Docker struct {
	Client      *client.Client
	Config      Config
	ContainerId string
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}
