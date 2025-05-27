package db

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// schema — SQL schema for the DB
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(255) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init opens the DB and creates a table if it doesn't exist
func Init(dbFile string) error {
	var err error

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Check if the DB is available
	if err = db.Ping(); err != nil {
		return err
	}

	// Always try to create a table (if it doesn't exist)
	_, err = db.Exec(schema)
	if err != nil {
		return errors.New("error creating DB schema: " + err.Error())
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
