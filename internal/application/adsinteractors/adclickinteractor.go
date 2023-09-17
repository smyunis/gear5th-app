package adsinteractors

import (
	"context"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/adclick"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdClickInteractor struct {
	adClickRepository  adclick.AdClickRepository
	cacheStore         application.KeyValueStore
	digitalSignService application.DigitalSignatureService
	eventDispatcher    application.EventDispatcher
	logger             application.Logger
}

func NewAdClickInteractor(
	adClickRepository adclick.AdClickRepository,
	cacheStore application.KeyValueStore,
	digitalSignService application.DigitalSignatureService,
	eventDispatcher application.EventDispatcher,
	logger application.Logger) AdClickInteractor {
	return AdClickInteractor{
		adClickRepository,
		cacheStore,
		digitalSignService,
		eventDispatcher,
		logger,
	}
}

func (i *AdClickInteractor) OnClick(adPieceID shared.ID, siteID shared.ID, adSlotID shared.ID, publisherID shared.ID, token string) error {
	if !i.digitalSignService.Validate(token) {
		return application.ErrAuthorization
	}

	viewID, err := i.digitalSignService.GetMessage(token)
	if err != nil {
		return err
	}

	_, err = i.cacheStore.Get(ViewIDCacheKey(viewID))
	if err != nil {
		return err
	}

	a := adclick.NewAdClick(shared.ID(viewID), adPieceID, siteID, adSlotID, publisherID)

	err = i.adClickRepository.Save(context.Background(), a)
	if err != nil {
		i.logger.Error("adclick/save", err)
		return err
	}

	i.eventDispatcher.DispatchAsync(a.Events)

	return nil

}
