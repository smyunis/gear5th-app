package advertiser

import (
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type Advertiser struct {
	ID    shared.ID
	Name  string
	Email user.Email
}

func NewAdvertiser(name string, email user.Email) Advertiser {
	return Advertiser{
		ID:    shared.NewID(),
		Name:  name,
		Email: email,
	}
}

func (a *Advertiser) StartCampaign(name string, start, end time.Time) campaign.Campaign {
	return campaign.NewCampaign(name, a.ID, start, end)
}
