package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go1f/pkg/db"
)

// Web handler for /api/task
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id not specified"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "task not found"})
		return
	}

	writeJson(w, task)
}

// Web handler for /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "error parsing JSON"})
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "id not specified"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title not specified"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]any{}) // empty JSON
}

// Web handler for /api/task/done
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id not specified"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "task not found"})
		return
	}

	if task.Repeat == "" {
		// One-time task — delete
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, map[string]any{}) // {}
		return
	}

	// Periodic — calculate and update date
	nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": "Date calculation error: " + err.Error()})
		return
	}

	if err := db.UpdateDate(nextDate, id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]any{})
}

// Web handler for /api/task method delete
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id parameter is missing"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, struct{}{})
}
