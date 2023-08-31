package siteinteractors_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/siteinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var interactor siteinteractors.SiteInteractor

func setup() {
	siteRepository := testdoubles.SiteRepositoryStub{}
	userRepo := testdoubles.UserRepositoryStub{}
	logger := testdoubles.ConsoleLogger{}
	ed := testdoubles.LocalizedEventDispatcher{}
	siteVerifService := testdoubles.AdsTxtSiteVerificaitonStub{}
	interactor = siteinteractors.NewSiteInteractor(siteRepository, userRepo, siteVerifService, &ed, logger)
}

func TestSiteInteractor(t *testing.T) {
	// Setup correct
}

func TestCreateNewSite(t *testing.T) {
	siteURL, _ := url.Parse("https://www.google.com")
	err := interactor.CreateSite(testdoubles.StubID, *siteURL)
	if err != nil {
		t.FailNow()
	}
}

func TestVerifySite(t *testing.T) {
	err := interactor.VerifySite(testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}
}

func TestGenerateAdsTxtRecord(t *testing.T) {
	record, err := interactor.GenerateAdsTxtRecord(testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}

	expectedRecord := fmt.Sprintf("gear5th.com, %s, DIRECT", testdoubles.StubID)

	if record.String() != expectedRecord {
		t.FailNow()
	}

}

func TestRemoveSite(t *testing.T) {
	err := interactor.RemoveSite(testdoubles.StubID, testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}
}

func TestGetActiveSitesForPublisher(t *testing.T) {
	_, err := interactor.ActiveSitesForPublisher(testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}
}
