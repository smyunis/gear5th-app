//go:build db
package adslotrepository_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/adslotpersistence/adslotrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var adslotRepo adslot.AdSlotRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	adslotRepo = adslotrepository.NewMongoDBAdSlotRepository(dbStore, logger)
}

func TestSaveAdslot(t *testing.T) {
	slot := adslot.NewAdSlot(shared.NewID(), "ad-slot-test-x", adslot.Vertical)
	err := adslotRepo.Save(context.TODO(), slot)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedAdSlot(t *testing.T) {
	slot := adslot.NewAdSlot(shared.NewID(), "ad-slot-test-xx", adslot.Box)
	err := adslotRepo.Save(context.TODO(), slot)
	if err != nil {
		t.Fatal(err)
	}

	fetchedAdSlot, err := adslotRepo.Get(context.TODO(), slot.ID())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fetchedAdSlot, slot) {
		t.Fatal(fetchedAdSlot, slot)
	}
}

func TestGetActiveAdSlotsForSite(t *testing.T) {
	siteID := shared.NewID()
	slot1 := adslot.NewAdSlot(siteID, "ad-slot-test-xx", adslot.Box)
	slot2 := adslot.NewAdSlot(siteID, "ad-slot-test-yy", adslot.Vertical)

	err := adslotRepo.Save(context.TODO(), slot1)
	if err != nil {
		t.Fatal(err)
	}
	err = adslotRepo.Save(context.TODO(), slot2)
	if err != nil {
		t.Fatal(err)
	}

	slots, err := adslotRepo.ActiveAdSlotsForSite(siteID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(slots,[]adslot.AdSlot{slot1,slot2}) {
		t.FailNow()
	}
}
