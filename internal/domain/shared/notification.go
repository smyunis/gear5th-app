package shared

import "time"

type Notification struct {
	id Id
	message string
	time time.Time
}

