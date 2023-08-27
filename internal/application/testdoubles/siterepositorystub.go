package testdoubles

import (
	"context"
	"net/url"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

const StubID = "stub-id-xxx"

type SiteRepositoryStub struct{}

// ActiveSitesForPublisher implements site.SiteRepository.
func (SiteRepositoryStub) ActiveSitesForPublisher(publisherID shared.ID) ([]site.Site, error) {
	u, _ := url.Parse("https://www.google.com")
	s1 := site.ReconstituteSite(StubID, publisherID, *u, true, false, []site.MonetizationStatus{{true, time.Now()}})
	u2, _ := url.Parse("https://www.bing.com")
	s2 := site.ReconstituteSite(StubID, publisherID, *u2, true, false, []site.MonetizationStatus{{true, time.Now()}})
	return []site.Site{s1, s2}, nil
}

func (SiteRepositoryStub) Get(ctx context.Context, id shared.ID) (site.Site, error) {
	u, _ := url.Parse("https://www.google.com")
	return site.ReconstituteSite(StubID, StubID, *u, true, false, []site.MonetizationStatus{{true, time.Now()}}), nil
}

func (SiteRepositoryStub) Save(ctx context.Context, e site.Site) error {
	return nil
}
