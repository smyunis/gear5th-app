//go:build db
package depositrepository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/deposit"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/publisherpersistence/depositrepository"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var depRepo deposit.DepositRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	depRepo = depositrepository.NewMongoDBDepositRepository(dbStore, logger)
}

func TestSaveDeposit(t *testing.T) {
	st := time.Now()
	ed := st.Add(360 * time.Hour)
	d := deposit.NewDeposit(shared.NewID(), 4566.5, st, ed)
	err := depRepo.Save(context.TODO(), d)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedDeposit(t *testing.T) {
	st := time.Now()
	ed := st.Add(363 * time.Hour)
	d := deposit.NewDeposit(shared.NewID(), 4566.5, st, ed)
	err := depRepo.Save(context.TODO(), d)
	if err != nil {
		t.Fatal(err)
	}

	fd, err := depRepo.Get(context.TODO(), d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if fd.ID != d.ID || fd.AdvertiserID != d.AdvertiserID {
		t.FailNow()
	}
}

func TestDepositsForToday(t *testing.T) {
	st := time.Now()
	ed := st.Add(363 * time.Hour)
	d := deposit.NewDeposit(shared.NewID(), 22334, st, ed)
	err := depRepo.Save(context.TODO(), d)
	if err != nil {
		t.Fatal(err)
	}

	dd,err := depRepo.DailyDisposits(st)
	if err != nil {
		t.Fatal(err)
	}

	if len(dd) == 0 {
		t.FailNow()
	}
}

