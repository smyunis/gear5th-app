package siteinteractors

import (
	"context"
	"errors"
	"net/url"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type SiteInteractor struct {
	siteRepository          site.SiteRepository
	siteVerificationService site.AdsTxtVerificationService
	eventDispatcher         application.EventDispatcher
	logger                  application.Logger
}

func NewSiteInteractor(siteRepository site.SiteRepository,
	siteVerificationService site.AdsTxtVerificationService,
	eventDispatcher application.EventDispatcher,
	logger application.Logger) SiteInteractor {
	return SiteInteractor{
		siteRepository,
		siteVerificationService,
		eventDispatcher,
		logger,
	}
}

func (i *SiteInteractor) CreateSite(publisherID shared.ID, siteURL url.URL) error {
	s := site.NewSite(publisherID, siteURL)
	defer i.eventDispatcher.DispatchAsync(s.DomainEvents())
	err := i.siteRepository.Save(context.Background(), s)
	if err != nil {
		i.logger.Error("site/createsite", err)
		return err
	}

	return nil
}

func (i *SiteInteractor) VerifySite(siteID shared.ID) error {
	s, err := i.siteRepository.Get(context.Background(), siteID)
	defer i.eventDispatcher.DispatchAsync(s.DomainEvents())

	if err != nil {
		if !errors.Is(err, application.ErrEntityNotFound) {
			i.logger.Error("site/verifysite", err)
		}
		return err
	}
	adsTxtRecord := site.GetAdsTxtRecord(s)
	err = i.siteVerificationService.VerifyAdsTxt(&s, adsTxtRecord)
	if err != nil {
		return site.ErrSiteVerification
	}
	i.siteRepository.Save(context.Background(), s)
	return nil
}
