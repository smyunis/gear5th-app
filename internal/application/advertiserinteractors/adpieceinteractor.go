package advertiserinteractors

import (
	"context"
	"net/url"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/adpiece"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/authorization"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdPieceInteractor struct {
	adPieceRepository  adpiece.AdPieceRepository
	campaignRepository campaign.CampaignRepository
	userRepository     user.UserRepository
	eventDispatcher    application.EventDispatcher
}

func NewAdPieceInteractor(adPieceRepository adpiece.AdPieceRepository,
	campaignRepository campaign.CampaignRepository,
	userRepository user.UserRepository,
	eventDispatcher application.EventDispatcher) AdPieceInteractor {
	return AdPieceInteractor{
		adPieceRepository,
		campaignRepository,
		userRepository,
		eventDispatcher,
	}
}

func (i *AdPieceInteractor) DeactivateAdPiece(actorID shared.ID, adPieceID shared.ID) error {

	actor, err := i.userRepository.Get(context.Background(), actorID)
	if err != nil {
		return err
	}

	a, err := i.adPieceRepository.Get(context.Background(), adPieceID)
	if err != nil {
		return err
	}

	c, err := i.campaignRepository.Get(context.Background(), a.CampaignID)
	if err != nil {
		return err
	}

	if !authorization.CanManageAdPiece(actor, c, a) {
		return application.ErrAuthorization
	}

	a.Deactivate()

	err = i.adPieceRepository.Save(context.Background(), a)
	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(a.Events)
	return nil
}

func (i *AdPieceInteractor) ResourceURL(adPieceID shared.ID) (*url.URL, error) {
	_, err := i.adPieceRepository.Get(context.Background(), adPieceID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *AdPieceInteractor) ActiveAdPiecesForCampaign(campaignID shared.ID) ([]adpiece.AdPiece, error) {
	return i.adPieceRepository.ActiveAdPiecesForCampaign(campaignID)
}

func (i *AdPieceInteractor) AdPiece(adPieceID shared.ID) (adpiece.AdPiece, error) {
	return i.adPieceRepository.Get(context.Background(), adPieceID)
}
