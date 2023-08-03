package publisher

import "gitlab.com/gear5th/gear5th-api/internal/domain/shared"

type Publisher struct {
	userId                      shared.Id
	acknowledgedNotifications   []shared.Notification
	unacknowledgedNotifications []shared.Notification
}

func NewPublisher(userId shared.Id) Publisher {
	return Publisher{
		userId: userId,
	}
}
