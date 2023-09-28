package shared

import "time"

func TimeEdges(start time.Time, end time.Time) (time.Time, time.Time) {
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	return s, e
}

func DayTimeEdges(day time.Time) (time.Time, time.Time) {
	s := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	e := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	return s, e
}


