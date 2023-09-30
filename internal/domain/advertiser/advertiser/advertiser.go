package advertiser

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdvertiserRepository interface {
	shared.EntityRepository[Advertiser]
	Advertisers() ([]Advertiser, error)
}

type Advertiser struct {
	UserID shared.ID
	Name   string
	Note   string
}

func NewAdvertiser(userID shared.ID, name string) Advertiser {
	return Advertiser{
		UserID: userID,
		Name:   name,
	}
}

func (a *Advertiser) StartCampaign(name string, start, end time.Time) campaign.Campaign {
	return campaign.NewCampaign(name, a.UserID, start, end)
}
