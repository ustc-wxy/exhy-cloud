package task

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	Cmd           []string
	Image         string
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy string
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

func NewDocker(config Config, containerId string) Docker {
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	return Docker{
		Client:      dc,
		Config:      config,
		ContainerId: containerId,
	}
}

// Run() performs the same operations as the `docker run` command
func (d *Docker) Run() DockerResult {
	ctx := context.Background()

	// Pull the specified image
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, types.ImagePullOptions{})
	if err != nil {
		log.Printf(
			"Error pulling image %s: %v\n",
			d.Config.Image, err,
		)
		return DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)

	// Pre-work for creating container
	rp := container.RestartPolicy{
		Name: d.Config.RestartPolicy,
	}

	r := container.Resources{
		Memory: d.Config.Memory,
	}

	containerConfig := container.Config{
		Image: d.Config.Image,
		Env:   d.Config.Env,
	}

	containerHostConfig := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	// Create the container
	resp, err := d.Client.ContainerCreate(
		ctx, &containerConfig, &containerHostConfig, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf(
			"Error creating container using image %s : %v\n",
			d.Config.Image, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(
		ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Printf(
			"Error starting container %s : %v",
			resp.ID, err)
		return DockerResult{Error: err}
	}

	d.ContainerId = resp.ID

	// Show container logs
	out, err := d.Client.ContainerLogs(
		ctx, resp.ID,
		types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("Error getting logs for container %s : %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		ContainerId: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

// Stop() performs the same operations as the `docker run` & `docker rm` command
func (d *Docker) Stop() DockerResult {
	ctx := context.Background()

	log.Printf(
		"Attempting to stop container %v",
		d.ContainerId,
	)

	// Stop container
	err := d.Client.ContainerStop(ctx, d.ContainerId, container.StopOptions{})
	if err != nil {
		panic(err)
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}

	// Remove container
	err = d.Client.ContainerRemove(
		ctx,
		d.ContainerId,
		removeOptions,
	)
	if err != nil {
		panic(err)
	}

	return DockerResult{
		Action: "stop",
		Result: "success",
		Error:  nil,
	}
}
