package worker

import (
	"encoding/json"
	"exhy-cloud/task"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}
	if err := d.Decode(&te); err != nil {
		log.Printf("[http] Error decode body: %v\n", err)
		w.WriteHeader(400)
		return
	}

	a.Worker.AddTask(te.Task)
	log.Printf("[http] Added task %v\n", te.Task.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
}

func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Println("[http] No taskID passed in request.")
		w.WriteHeader(400)
	}
	tID, _ := uuid.Parse(taskID)
	if _, ok := a.Worker.Db[tID]; !ok {
		log.Printf("[http] No task ID %v found.", tID)
		w.WriteHeader(404)
		return
	}

	taskEntry := a.Worker.Db[tID]
	taskCopy := *taskEntry
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)

	log.Printf("[http] Task %v on container %v will be stopped", taskEntry.ID, taskEntry.ContainerID)
	w.WriteHeader(204)
}
