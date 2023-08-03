package publisher

import "gitlab.com/gear5th/gear5th-api/internal/domain/shared"

type Publisher struct {
	userId                      shared.Id
	unacknowledgedNotifications []shared.Notification
}

func NewPublisher(userId shared.Id) Publisher {
	return Publisher{
		userId: userId,
	}
}

func (p *Publisher) Notify(notification shared.Notification) {
	p.unacknowledgedNotifications = append(p.unacknowledgedNotifications, notification)
}

func (p *Publisher) UnacknowledgedNotifications() []shared.Notification {
	notifications := make([]shared.Notification, len(p.unacknowledgedNotifications))
	copy(notifications, p.unacknowledgedNotifications)
	return notifications
}

func (p *Publisher) AcknowledgeNotifications() {
	p.unacknowledgedNotifications = nil
}
