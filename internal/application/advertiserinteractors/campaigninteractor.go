package advertiserinteractors

import (
	"context"
	"io"
	"net/url"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/adpiece"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/authorization"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type CampaignInteractor struct {
	campaignRepository campaign.CampaignRepository
	userRepository     user.UserRepository
	adPieceRepository  adpiece.AdPieceRepository
	fileStore          application.FileStore
	eventDispatcher    application.EventDispatcher
}

func NewCampaignInteractor(
	campaignRepository campaign.CampaignRepository,
	userRepository user.UserRepository,
	adPieceRepository adpiece.AdPieceRepository,
	fileStore application.FileStore,
	eventDispatcher application.EventDispatcher) CampaignInteractor {
	return CampaignInteractor{
		campaignRepository,
		userRepository,
		adPieceRepository,
		fileStore,
		eventDispatcher,
	}
}

func (i *CampaignInteractor) Campaign(campaignID shared.ID) (campaign.Campaign, error) {
	return i.campaignRepository.Get(context.Background(), campaignID)
}


func (i *CampaignInteractor) StartCampaign(actorID shared.ID, name string, start, end time.Time) error {

	actor, err := i.userRepository.Get(context.Background(), actorID)
	if err != nil {
		return err
	}

	if !authorization.CanStartCampaign(actor) {
		return application.ErrAuthorization
	}

	c := campaign.NewCampaign(name, actorID, start, end)
	err = i.campaignRepository.Save(context.Background(), c)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(c.Events)
	return nil
}

func (i *CampaignInteractor) QuitCampaign(actorID shared.ID, campaignID shared.ID) error {

	actor, err := i.userRepository.Get(context.Background(), actorID)
	if err != nil {
		return err
	}

	c, err := i.campaignRepository.Get(context.Background(), campaignID)
	if err != nil {
		return err
	}

	if !authorization.CanManageCampaign(actor, c) {
		return application.ErrAuthorization
	}

	c.Quit()

	err = i.campaignRepository.Save(context.Background(), c)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(c.Events)
	return nil
}

func (i *CampaignInteractor) AddAdPiece(actorID shared.ID, campaignID shared.ID,
	slot adslot.AdSlotType, ref *url.URL, resourceMIME string, resource io.Reader) error {

	actor, err := i.userRepository.Get(context.Background(), actorID)
	if err != nil {
		return err
	}

	c, err := i.campaignRepository.Get(context.Background(), campaignID)
	if err != nil {
		return err
	}

	if !authorization.CanManageCampaign(actor, c) {
		return application.ErrAuthorization
	}

	resourceID, err := i.fileStore.Save(resource)
	if err != nil {
		return err
	}

	adPiece := c.AddAdPiece(slot, ref, resourceID, resourceMIME)

	err = i.adPieceRepository.Save(context.Background(), adPiece)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(c.Events, adPiece.Events)

	return nil
}

func (i *CampaignInteractor) CampaignsForAdvertiser(advertiserID shared.ID) ([]campaign.Campaign, error) {
	return i.campaignRepository.CampaignsForAdvertiser(advertiserID)
}
