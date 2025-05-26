package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DateFormat — format for the date
const DateFormat = "20060102"

// Web handler for /api/nextdate
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	// if nowStr is empty, use current time
	if nowStr == "" {
		now = time.Now()
	} else {
		// parse nowStr to time
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid 'now' parameter", http.StatusBadRequest)
			return
		}
	}

	// calculate the next date
	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}

// afterNow — true, if date > now
func afterNow(now, t time.Time) bool {
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return now.After(t)
}

// NextDate — calculates the next date of the task by the rule
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repetition not set")
	}

	// parse the start date
	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("invalid date format dstart")
	}

	// parse the repeat rule
	parts := strings.Split(repeat, " ")
	switch parts[0] {

	// if the rule is 'd'
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid format of rule 'd'")
		}
		// parse the number of days
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid number of days")
		}
		// calculate the next date
		for {
			startDate = startDate.AddDate(0, 0, days)
			if afterNow(startDate, now) {
				break
			}

		}
		return startDate.Format(DateFormat), nil

	// if the rule is 'y'
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid format of rule 'y'")
		}
		// calculate the next date
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if afterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	// if the rule is not supported
	default:
		return "", errors.New("unsupported repetition rule")
	}
}
