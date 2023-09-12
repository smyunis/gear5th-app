package authorization

import (
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
)

func CanStartCampaign(actor user.User) bool {
	return actor.HasRole(user.Administrator) || actor.HasRole(user.Advertiser)
}

func CanManageCampaign(actor user.User, campaign campaign.Campaign) bool {
	if actor.HasRole(user.Advertiser) {
		return campaign.AdvertiserUserID == actor.UserID()
	}
	return actor.HasRole(user.Administrator)
}
