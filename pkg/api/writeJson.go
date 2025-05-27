package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON — helper function to return JSON response
func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
	}
}
