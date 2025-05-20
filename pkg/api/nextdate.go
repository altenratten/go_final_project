package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// afterNow — true, если date > now (без учёта времени)
func afterNow(now, t time.Time) bool {
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return now.After(t)
}

// NextDate — вычисляет следующую дату задачи по правилу
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("повторение не задано")
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("неверный формат даты dstart")
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("неверное значение дней")
		}
		for {
			startDate = startDate.AddDate(0, 0, days)
			if afterNow(startDate, now) {
				break
			}

		}
		return startDate.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат правила y")
		}
		for {
			startDate = startDate.AddDate(1, 0, 0)
			if afterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(DateFormat), nil

	default:
		return "", errors.New("неподдерживаемое правило повторения")
	}
}
