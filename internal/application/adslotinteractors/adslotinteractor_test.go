package adslotinteractors_test

import (
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/adslotinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var interactor adslotinteractors.AdSlotInteractor

func setup() {
	siteRepository := testdoubles.SiteRepositoryStub{}
	userRepo := testdoubles.UserRepositoryStub{}
	adslotrepo := testdoubles.AdSlotRepositoryStub{}
	ed := testdoubles.LocalizedEventDispatcher{}
	interactor = adslotinteractors.NewAdSlotInteractor(siteRepository, userRepo, adslotrepo, &ed)
}

func TestCreateAdSlot(t *testing.T) {
	err := interactor.CreateAdSlotForSite(testdoubles.StubID, testdoubles.StubID, "my-site-adslot", adslot.Box)
	if err != nil {
		t.FailNow()
	}
}

func TestChangeAdSlotName(t *testing.T) {
	err := interactor.ChangeAdSlotName(testdoubles.StubID, testdoubles.StubID, "new name")
	if err != nil {
		t.FailNow()
	}
}

func TestGetIntegrationHTMLSnippet(t *testing.T) {
	_, err := interactor.GetIntegrationHTMLSnippet(testdoubles.StubID)
	if err != nil {
		t.FailNow()
	}
}
