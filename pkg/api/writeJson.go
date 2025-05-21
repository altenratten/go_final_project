package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON — helper function to return JSON response
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
