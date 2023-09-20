//go:build db
package adclickrepository_test

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/ads/adclick"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/adspersistence/adclickrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var adClickRepo adclick.AdClickRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	adClickRepo = adclickrepository.NewMongoDBAdClickRepository(dbStore, logger)
}

func TestSaveAdClick(t *testing.T) {
	i := shared.NewID()
	a := adclick.NewAdClick(i, i, i, i, i)

	err := adClickRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedSlotId(t *testing.T) {
	i := shared.NewID()
	a := adclick.NewAdClick(i, i, i, i, i)

	err := adClickRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}

	fa, err := adClickRepo.Get(context.TODO(), a.ID)

	if fa.ID != a.ID || fa.ViewID != a.ViewID {
		t.FailNow()
	}
}

func TestAdclicksWithintimeperiod(t *testing.T) {
	i := shared.NewID()
	a := adclick.NewAdClick(i, i, i, i, i)

	err := adClickRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}
	a2 := adclick.NewAdClick(i, i, i, i, i)

	err = adClickRepo.Save(context.TODO(), a2)
	if err != nil {
		t.Fatal(err)
	}
	s := time.Date(2022, 1, 1, 1, 1, 1, 1, time.Local)
	adc, err := adClickRepo.AdClicksForPublisher(i, s, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	if !slices.ContainsFunc(adc, func(ac adclick.AdClick) bool {
		return ac.ViewID == i
	}) {
		t.FailNow()
	}
}
