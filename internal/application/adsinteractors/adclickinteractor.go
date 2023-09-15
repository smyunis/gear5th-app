package adsinteractors

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdClickInteractor struct {
	cacheStore         application.KeyValueStore
	digitalSignService application.DigitalSignatureService
	logger             application.Logger
}

func NewAdClickInteractor(
	cacheStore application.KeyValueStore,
	digitalSignService application.DigitalSignatureService,
	logger application.Logger) AdClickInteractor {
	return AdClickInteractor{
		cacheStore,
		digitalSignService,
		logger,
	}
}

func (i *AdClickInteractor) OnClick(adPieceID shared.ID, token string) error {
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

	// save to db here
	// emit event for payment

	return nil

}
