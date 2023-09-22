//go:build db
package earningrepository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/earning"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/paymentpersistence/earningrepository"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var earningRepo earning.EarningRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	earningRepo = earningrepository.NewMongoDBEarningRepository(dbStore, logger)
}

func TestSaveEarning(t *testing.T) {
	i := shared.NewID()
	e := earning.NewEarning(i, earning.Impression, 5000, i, i, i)
	err := earningRepo.Save(context.TODO(), e)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedEarning(t *testing.T) {
	i := shared.NewID()
	e := earning.NewEarning(i, earning.Impression, 6000, i, i, i)
	err := earningRepo.Save(context.TODO(), e)
	if err != nil {
		t.Fatal(err)
	}

	fe, err := earningRepo.Get(context.TODO(), e.ID)
	if err != nil {
		t.Fatal(err)
	}

	if fe.AdPieceID != e.AdPieceID {
		t.FailNow()
	}
}

func TestPublisherEarnings(t *testing.T) {
	i := shared.NewID()
	ti := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
	e := earning.NewEarning(i, earning.Impression, 6000, i, i, i)
	err := earningRepo.Save(context.TODO(), e)
	if err != nil {
		t.Fatal(err)
	}

	e2 := earning.NewEarning(i, earning.Impression, 8000, i, i, i)
	err = earningRepo.Save(context.TODO(), e2)
	if err != nil {
		t.Fatal(err)
	}

	pe, err := earningRepo.EarningsForPublisher(i, ti, time.Now().Add(500*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	if len(pe) == 0 {
		t.FailNow()
	}
}
