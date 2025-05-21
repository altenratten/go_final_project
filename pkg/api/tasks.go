package api

import (
	"go1f/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Web handler for /api/tasks
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}
	writeJson(w, struct {
		Tasks []*db.Task `json:"tasks"`
	}{Tasks: tasks})
}
