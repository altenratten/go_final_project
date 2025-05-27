package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// DateFormat — format for the date
const DateFormat = "20060102"
const baseLimit = 50

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask — adds a task to the DB
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// parseSearchDate — parses the date from the string
func parseSearchDate(input string) (string, error) {
	t, err := time.Parse("02.01.2006", input)
	if err != nil {
		return "", err
	}
	return t.Format(DateFormat), nil
}

// GetTask — gets a task from the DB
func GetTask(id string) (*Task, error) {
	var t Task
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTask — updates a task in the DB
func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

// DeleteTask — deletes a task from the DB
func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("task not found")
	}
	return nil
}

// UpdateDate — updates the date of a task in the DB
func UpdateDate(next string, id string) error {
	_, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	return err
}
