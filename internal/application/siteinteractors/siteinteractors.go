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
	err := i.siteRepository.Save(context.Background(), s)
	if err != nil {
		i.logger.Error("site/createsite", err)
		return err
	}
	i.eventDispatcher.DispatchAsync(s.DomainEvents())

	return nil
}

func (i *SiteInteractor) VerifySite(siteID shared.ID) error {
	s, err := i.siteRepository.Get(context.Background(), siteID)

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
	i.eventDispatcher.DispatchAsync(s.DomainEvents())

	return nil
}

func (i *SiteInteractor) GenerateAdsTxtRecord(siteID shared.ID) (site.AdsTxtRecord, error) {
	s, err := i.siteRepository.Get(context.Background(), siteID)
	if err != nil {
		return site.AdsTxtRecord{}, err
	}
	record := site.GetAdsTxtRecord(s)
	return record, nil
}

func (i *SiteInteractor) RemoveSite(siteID shared.ID) error {
	s, err := i.siteRepository.Get(context.Background(), siteID)
	if err != nil {
		return err
	}
	s.Deactivate()
	i.siteRepository.Save(context.Background(), s)
	i.eventDispatcher.DispatchAsync(s.DomainEvents())
	return nil
}

func (i *SiteInteractor) ActiveSitesForPublisher(publisherID shared.ID) ([]site.Site, error) {
	return i.siteRepository.ActiveSitesForPublisher(publisherID)
}
