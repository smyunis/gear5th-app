package adsinteractors

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/adclick"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/impression"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdsInteractor struct {
	impressionRepository impression.ImpressionRepository
	adClickRepository    adclick.AdClickRepository
	siteRepository       site.SiteRepository
	adSlotRepository     adslot.AdSlotRepository
	cacheStore           application.KeyValueStore
	digitalSignService   application.DigitalSignatureService
	eventDispatcher      application.EventDispatcher
	logger               application.Logger
}

func NewAdsInteractor(
	impressionRepository impression.ImpressionRepository,
	adClickRepository adclick.AdClickRepository,
	siteRepository site.SiteRepository,
	adSlotRepository adslot.AdSlotRepository,
	cacheStore application.KeyValueStore,
	digitalSignService application.DigitalSignatureService,
	eventDispatcher application.EventDispatcher,
	logger application.Logger) AdsInteractor {
	return AdsInteractor{
		impressionRepository,
		adClickRepository,
		siteRepository,
		adSlotRepository,
		cacheStore,
		digitalSignService,
		eventDispatcher,
		logger,
	}
}

func (i *AdsInteractor) OnImpression(adPieceID shared.ID, siteID shared.ID, adSlotID shared.ID, publisherID shared.ID, token string) error {

	viewID, err := i.validateToken(token)
	if err != nil {
		return application.ErrRequirementFailed
	}

	if !i.siteCanServeAdPieces(siteID) {
		return application.ErrRequirementFailed
	}

	if !i.adSlotCanServeAdPieces(adSlotID) {
		return application.ErrRequirementFailed
	}

	a := impression.NewImpression(shared.ID(viewID), adPieceID, siteID, adSlotID, publisherID)

	err = i.impressionRepository.Save(context.Background(), a)
	if err != nil {
		i.logger.Error("impression/save", err)
		return err
	}

	i.eventDispatcher.DispatchAsync(a.Events)

	return nil
}

func (i *AdsInteractor) OnClick(adPieceID shared.ID, siteID shared.ID, adSlotID shared.ID, publisherID shared.ID, token string) error {
	viewID, err := i.validateToken(token)
	if err != nil {
		return application.ErrRequirementFailed
	}

	if !i.siteCanServeAdPieces(siteID) {
		return application.ErrRequirementFailed
	}

	if !i.adSlotCanServeAdPieces(adSlotID) {
		return application.ErrRequirementFailed
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

func (i *AdsInteractor) validateToken(token string) (string, error) {
	if !i.digitalSignService.Validate(token) {
		return "", application.ErrAuthorization
	}

	viewID, err := i.digitalSignService.GetMessage(token)
	if err != nil {
		return "", err
	}

	_, err = i.cacheStore.Get(ViewIDCacheKey(viewID))
	if err != nil {
		return "", err
	}
	return viewID, nil
}

func (i *AdsInteractor) siteCanServeAdPieces(siteID shared.ID) bool {
	siteCanServeStatusCacheKey := fmt.Sprintf("site:%s:canServeAdPiece", siteID.String())
	sv, err := i.cacheStore.Get(siteCanServeStatusCacheKey)
	canServe, parseErr := strconv.ParseBool(sv)
	if err != nil || parseErr != nil {
		s, err := i.siteRepository.Get(context.Background(), siteID)
		if err != nil {
			i.logger.Error("impression/site/get", err)
			return false
		}

		i.cacheStore.Save(siteCanServeStatusCacheKey, strconv.FormatBool(s.CanServeAdPiece()), 3*time.Hour)
		return s.CanServeAdPiece()
	}
	return canServe
}

func (i *AdsInteractor) adSlotCanServeAdPieces(slotID shared.ID) bool {
	adSlotCanServeCacheKey := fmt.Sprintf("adslot:%s:canServeAdPiece", slotID.String())
	sv, err := i.cacheStore.Get(adSlotCanServeCacheKey)
	canServe, parseErr := strconv.ParseBool(sv)
	if err != nil || parseErr != nil {
		s, err := i.adSlotRepository.Get(context.Background(), slotID)
		if err != nil {
			i.logger.Error("impression/adslot/get", err)
			return false
		}

		i.cacheStore.Save(adSlotCanServeCacheKey, strconv.FormatBool(s.CanServeAdPieces()), 3*time.Hour)
		return s.CanServeAdPieces()
	}
	return canServe
}
