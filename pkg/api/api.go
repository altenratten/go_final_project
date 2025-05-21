package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// handling POST method
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
