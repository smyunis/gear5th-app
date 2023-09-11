package advertiserinteractors

import (
	"context"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type CampaignInteractor struct {
	campaignRepository campaign.CampaignRepository
	eventDispatcher    application.EventDispatcher
}

func NewCampaignInteractor(
	campaignRepository campaign.CampaignRepository,
	eventDispatcher application.EventDispatcher) CampaignInteractor {
	return CampaignInteractor{
		campaignRepository,
		eventDispatcher,
	}
}

func (i *CampaignInteractor) StartCampaign(advertiserID shared.ID, name string, start, end time.Time) {
	c := campaign.NewCampaign(name, advertiserID, start, end)

	i.campaignRepository.Save(context.Background(), c)

	i.eventDispatcher.DispatchAsync(c.Events)
}
