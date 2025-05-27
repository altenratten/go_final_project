package main

import (
	"log"
	"os"

	"go1f/pkg/db"
	"go1f/pkg/server"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

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
	// Ensure the DB is closed when the application exits
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing the DB: %v", err)
		}
	}()

	// Start the web server
	server.StartServer()
}
