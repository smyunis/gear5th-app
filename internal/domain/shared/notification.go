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

func ReconstituteNotification(message string, t time.Time) Notification {
	return Notification{
		message,
		t,
	}
}

func (n Notification) Message() string {
	return n.message
}

func (n Notification) Time() time.Time {
	return n.time
}
