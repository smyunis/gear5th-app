package advertiser

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type Advertiser struct {
	UserID    shared.ID
	Name  string
}

func NewAdvertiser(name string) Advertiser {
	return Advertiser{
		UserID:    shared.NewID(),
		Name:  name,
	}
}

func (a *Advertiser) StartCampaign(name string, start, end time.Time) campaign.Campaign {
	return campaign.NewCampaign(name, a.UserID, start, end)
}
