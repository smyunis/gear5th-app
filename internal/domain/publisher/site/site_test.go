package site_test

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

func TestSiteHostnameVerification(t *testing.T) {
	siteUrl, _ := url.Parse("https://dev.to/tooljet/build-an-aws-s3-browser-with-tooljet-56d4")
	s := site.NewSite(shared.NewID(), *siteUrl)
	source, _ := url.Parse("https://dev.to/tooljet/build-an-aws-s3-browser-with-tooljet-56d4")

	err := site.VerifySiteHostname(&s, *source)

	if err != nil {
		t.FailNow()
	}
}

func TestGetAdsTxtRecord(t *testing.T) {
	siteUrl, _ := url.Parse("https://dev.to/tooljet/build-an-aws-s3-browser-with-tooljet-56d4")
	publisherId := shared.NewID()
	s := site.NewSite(publisherId, *siteUrl)
	expectedRecord := fmt.Sprintf("gear5th.com, %s, DIRECT", publisherId.String())

	record := site.GetAdsTxtRecord(s).String()

	if expectedRecord != record {
		t.Fatalf("expected: %s, given: %s", expectedRecord, record)
	}
}

func TestAddAdSlot(t *testing.T) {
	s := newSite()
	adslotName := "right-banner-big"
	var adslotType adslot.AdSlotType = adslot.Horizontal
	slot := s.AddAdSlot(adslotName, adslotType)

	if slot.Name != adslotName {
		t.FailNow()
	}
}

func TestDeactivateSite(t *testing.T) {
	s := newSite()

	s.Deactivate()

	if !s.IsDeactivated {
		t.FailNow()
	}
}

func TestCanNotServeAdPieceIfDeactvated(t *testing.T) {
	s := newSite()
	s.Deactivate()

	if s.CanServeAdPiece() {
		t.FailNow()
	}
}

func TestNotCanMonetizeIfDemonetized(t *testing.T) {
	s := newSite()
	s.Demonetize()

	if s.CanMonetize() {
		t.FailNow()
	}
}

func TestCanMonetizeIfMonetizationAllowedWhenCreation(t *testing.T) {
	s := newSite()
	s.Verify()

	if !s.CanMonetize() {
		t.FailNow()
	}
}

func TestCanNotDemonetizeConsecutively(t *testing.T) {
	s := newSite()

	s.Demonetize()
	s.Demonetize()

	history := s.MonetizationStatusHistory()
	if len(history) != 2 || history[len(history)-1].CanMonetize {
		t.FailNow()
	}
}

func TestCanNotAllowMonetizationStatusConsecutively(t *testing.T) {
	s := newSite()
	// New site already starts with monietization allowed
	s.AllowMonetization()
	s.AllowMonetization()

	history := s.MonetizationStatusHistory()
	if len(history) != 1 || !history[len(history)-1].CanMonetize {
		t.FailNow()
	}
}

func TestDemonetizeForTimePeriod(t *testing.T) {
	s := newSite()
	s.Verify()
	var demonetizationPeriod time.Duration = 2 * time.Second
	s.DemonetizeForTimePeriod(demonetizationPeriod)
	if s.CanMonetize() {
		t.FailNow()
	}
	time.Sleep(3 * time.Second)
	if !s.CanMonetize() {
		t.FailNow()
	}
}

func TestGetSiteDomain(t *testing.T) {
	s := newSite()
	domain := s.SiteDomain()
	if domain != "dev.to" {
		t.Fatal(domain)
	}

}

func newSite() site.Site {
	siteUrl, _ := url.Parse("https://dev.to/tooljet/build-an-aws-s3-browser-with-tooljet-56d4")
	publisherId := shared.NewID()
	return site.NewSite(publisherId, *siteUrl)
}
