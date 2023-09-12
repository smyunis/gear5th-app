package publisherinteractors_test

import (
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
)

var adSlotInteractor publisherinteractors.AdSlotInteractor

func adSlotTestSetup() {
	siteRepository := testdoubles.SiteRepositoryStub{}
	userRepo := testdoubles.UserRepositoryStub{}
	adslotrepo := testdoubles.AdSlotRepositoryStub{}
	ed := testdoubles.LocalizedEventDispatcher{}
	adSlotInteractor = publisherinteractors.NewAdSlotInteractor(siteRepository, userRepo, adslotrepo, &ed)
}

func TestCreateAdSlot(t *testing.T) {
	err := adSlotInteractor.CreateAdSlotForSite(testdoubles.StubID, testdoubles.StubID, "my-site-adslot", adslot.Box)
	if err != nil {
		t.FailNow()
	}
}

func TestChangeAdSlotName(t *testing.T) {
	err := adSlotInteractor.ChangeAdSlotName(testdoubles.StubID, testdoubles.StubID, "new name")
	if err != nil {
		t.FailNow()
	}
}

func TestGetIntegrationHTMLSnippet(t *testing.T) {
	_, err := adSlotInteractor.GetIntegrationHTMLSnippet(testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}
}
