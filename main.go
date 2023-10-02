package main

import (
	"exhy-cloud/task"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker/client"
)

func main() {
	fmt.Println("[Main] start...")
	d, r := createContainer()
	if r.Error != nil {
		log.Printf("[Main] Create container failed: %v\n", r.Error)
		os.Exit(1)
	}
	time.Sleep(time.Second * 5)
	log.Printf("[Main] stopping container %s\n", r.ContainerId)
	r = stopContainer(d)
	if r.Error != nil {
		log.Printf("[Main] Stop container failed: %v\n", r.Error)
		os.Exit(1)
	}
}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "mysql:5.7",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=123456",
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("create container failed: %v\n", result.Error)
		return nil, nil
	}

	fmt.Printf(
		"Container %s is running with config %v\n",
		result.ContainerId, c)
	return &d, &result
}

func stopContainer(d *task.Docker) *task.DockerResult {
	res := d.Stop()
	if res.Error != nil {
		return nil
	}
	fmt.Printf(
		"Container %s has been stopped and removed\n",
		res.ContainerId,
	)
	return &res
}
