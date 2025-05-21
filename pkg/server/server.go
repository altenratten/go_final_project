package server

import (
	"go1f/pkg/api"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	// Default port value
	port := "7540"

	// If the TODO_PORT environment variable is set, use it
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	api.Init()

	// File server for the web directory
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Printf("Server started on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
