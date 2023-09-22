package publisher

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type PublisherRepository interface {
	shared.EntityRepository[Publisher]
}

type Publisher struct {
	UserID           shared.ID
	LastDisbursement time.Time
}

func NewPublisher(userId shared.ID) Publisher {
	return Publisher{
		UserID:           userId,
		LastDisbursement: time.Now(),
	}
}

func (p *Publisher) UpdateLastDisbursement() {
	p.LastDisbursement = time.Now()
}
