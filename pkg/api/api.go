package api

import (
	"net/http"
)

// Init — initializes the API
func Init() {
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)
}

// Web handler for /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// GET method
	case http.MethodGet:
		getTaskHandler(w, r)
	// PUT method
	case http.MethodPut:
		updateTaskHandler(w, r)
	// POST method
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	// method not allowed
	default:
		writeJson(w, map[string]string{"error": "method not allowed"})
	}
}
