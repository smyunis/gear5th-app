package shared

import (
	"time"
)

type Notification struct {
	message string
	time    time.Time
}

func NewNotification(message string) Notification {
	return Notification{
		time:    time.Now(),
		message: message,
	}
}
