package db

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

var db *sql.DB

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

// Init открывает БД и создаёт таблицу, если её нет
func Init(dbFile string) error {
	var err error

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Проверяем, доступна ли БД
	if err = db.Ping(); err != nil {
		return err
	}

	// Всегда пытаемся создать таблицу (если её нет)
	_, err = db.Exec(schema)
	if err != nil {
		return errors.New("ошибка создания схемы БД: " + err.Error())
	}

	return nil
}
