package adclick

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdClickRepository interface {
	shared.EntityRepository[AdClick]
	AdClicksForPublisher(publisherID shared.ID, start time.Time, end time.Time) ([]AdClick, error)
	AdClicksCountForPublisher(publisherID shared.ID, start time.Time, end time.Time) (int, error)
}

type AdClick struct {
	ID                shared.ID
	Events            shared.Events
	ClickedAdPieceID  shared.ID
	OriginSiteID      shared.ID
	OriginAdSlotID    shared.ID
	OriginPublisherID shared.ID
	Time              time.Time
}

func NewAdClick(viewID shared.ID, clickedAdPieceID shared.ID, originSiteID shared.ID,
	originAdSlotID shared.ID,
	originPublisherID shared.ID) AdClick {
	a := AdClick{
		ID:                viewID,
		ClickedAdPieceID:  clickedAdPieceID,
		OriginSiteID:      originSiteID,
		OriginAdSlotID:    originAdSlotID,
		OriginPublisherID: originPublisherID,
		Time:              time.Now(),
		Events:            make(shared.Events),
	}
	a.Events.Emit("ad/clicked", a)
	return a
}
