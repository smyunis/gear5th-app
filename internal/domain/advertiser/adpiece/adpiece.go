package adpiece

import (
	"net/url"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdPieceRepository interface {
	shared.EntityRepository[AdPiece]
	ActiveAdPiecesForCampaign(campaignID shared.ID) ([]AdPiece, error)
}

type AdPiece struct {
	ID            shared.ID
	CampaignID    shared.ID
	SlotType      adslot.AdSlotType
	Resource      string
	Description   string
	Ref           *url.URL
	IsDeactivated bool
}

func NewAdPiece(campaignID shared.ID, slot adslot.AdSlotType, ref *url.URL, desc, resource string) AdPiece {
	return AdPiece{
		ID:          shared.NewID(),
		CampaignID:  campaignID,
		SlotType:    slot,
		Resource:    resource,
		Ref:         ref,
		Description: desc,
	}
}
