package api

import (
	"net/http"
)

// Init — initializes the API
func Init() {
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
	http.HandleFunc("/api/tasks", auth(getTasksHandler))
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signinHandler)
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
	//DELETE method
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	// method not allowed
	default:
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}
