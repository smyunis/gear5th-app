package impression

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type ImpressionRepository interface {
	shared.EntityRepository[Impression]
	ImpressionsForPublisher(publisherID shared.ID, start time.Time, end time.Time) ([]Impression, error)
}

type Impression struct {
	ID                shared.ID
	Events            shared.Events
	AdPieceID         shared.ID
	OriginSiteID      shared.ID
	OriginAdSlotID    shared.ID
	OriginPublisherID shared.ID
	Time              time.Time
}

func NewImpression(viewID shared.ID, adPieceID shared.ID, originSiteID shared.ID,
	originAdSlotID shared.ID,
	originPublisherID shared.ID) Impression {
	a := Impression{
		ID:                viewID,
		AdPieceID:         adPieceID,
		OriginSiteID:      originSiteID,
		OriginAdSlotID:    originAdSlotID,
		OriginPublisherID: originPublisherID,
		Time:              time.Now(),
		Events:            make(shared.Events),
	}
	a.Events.Emit("ad/impression-made", a)
	return a
}
