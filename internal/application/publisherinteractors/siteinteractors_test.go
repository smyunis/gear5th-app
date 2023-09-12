package publisherinteractors_test

import (
	"fmt"
	"net/url"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
)



var siteInteractor publisherinteractors.SiteInteractor

func siteTestSetup() {
	siteRepository := testdoubles.SiteRepositoryStub{}
	userRepo := testdoubles.UserRepositoryStub{}
	logger := testdoubles.ConsoleLogger{}
	ed := testdoubles.LocalizedEventDispatcher{}
	siteVerifService := testdoubles.AdsTxtSiteVerificaitonStub{}
	siteInteractor = publisherinteractors.NewSiteInteractor(siteRepository, userRepo, siteVerifService, &ed, logger)
}

func TestSiteInteractor(t *testing.T) {
	// Setup correct
}

func TestCreateNewSite(t *testing.T) {
	siteURL, _ := url.Parse("https://www.google.com")
	err := siteInteractor.CreateSite(testdoubles.StubID, *siteURL)
	if err != nil {
		t.FailNow()
	}
}

func TestVerifySite(t *testing.T) {
	err := siteInteractor.VerifySite(testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}
}

func TestGenerateAdsTxtRecord(t *testing.T) {
	record, err := siteInteractor.GenerateAdsTxtRecord(testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}

	expectedRecord := fmt.Sprintf("gear5th.com, %s, DIRECT", testdoubles.StubID)

	if record.String() != expectedRecord {
		t.FailNow()
	}

}

func TestRemoveSite(t *testing.T) {
	err := siteInteractor.RemoveSite(testdoubles.StubID, testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}
}

func TestGetActiveSitesForPublisher(t *testing.T) {
	_, err := siteInteractor.ActiveSitesForPublisher(testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}
}
