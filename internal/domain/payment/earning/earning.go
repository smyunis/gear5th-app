package earning

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type EarningRepository interface {
	shared.EntityRepository[Earning]
	EarningsForPublisher(publisherID shared.ID, start time.Time, end time.Time) ([]Earning, error)
}

type EarningReason int

const (
	_ EarningReason = iota
	Impression
	AdClick
)

type Earning struct {
	ID          shared.ID
	Events      shared.Events
	PublisherID shared.ID
	Reason      EarningReason
	AdPieceID   shared.ID
	AdSlotID    shared.ID
	SiteID      shared.ID
	Time        time.Time
	Amount      float64
}

func NewEarning(pubID shared.ID, reason EarningReason, amount float64, adPieceID, slotID, siteID shared.ID) Earning {
	return Earning{
		ID:          shared.NewID(),
		Events:      make(shared.Events),
		PublisherID: pubID,
		Reason:      reason,
		AdPieceID:   adPieceID,
		AdSlotID:    slotID,
		SiteID:      siteID,
		Time:        time.Now(),
		Amount:      amount,
	}
}
