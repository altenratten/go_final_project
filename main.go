package main

import (
	"log"
	"os"

	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	// Get database path from TODO_DBFILE environment variable
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Initialize the DB
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Error initializing the DB: %v", err)
	}

	// Start the web server
	server.StartServer()
}
