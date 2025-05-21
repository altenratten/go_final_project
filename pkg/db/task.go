package db

import (
	"database/sql"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

// Tasks — gets tasks from the DB
func Tasks(limit int, search string) ([]*Task, error) {
	const baseLimit = 50
	if limit <= 0 || limit > baseLimit {
		limit = baseLimit
	}

	var rows *sql.Rows
	var err error

	// if search is empty, get all tasks
	if search == "" {
		rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", limit)
	} else if date, errParse := parseSearchDate(search); errParse == nil {
		// search by date
		rows, err = db.Query("SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?", date, limit)
	} else {
		// search by title or comment
		like := "%" + search + "%"
		rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?", like, like, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

// parseSearchDate — parses the date from the string
func parseSearchDate(input string) (string, error) {
	t, err := time.Parse("02.01.2006", input)
	if err != nil {
		return "", err
	}
	return t.Format("20060102"), nil
}
