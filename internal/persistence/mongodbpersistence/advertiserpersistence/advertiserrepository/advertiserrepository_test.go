/// go:build db
package advertiserrepository_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/advertiser"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/advertiserpersistence/advertiserrepository"
	"gitlab.com/gear5th/gear5th-app/internal/persistence/mongodbpersistence/mongotestdoubles"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

var advertiserRepo advertiser.AdvertiserRepository

func setup() {
	configProvider := mongotestdoubles.NewTestEnvConfigurationProvider()
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	logger := testdoubles.ConsoleLogger{}
	advertiserRepo = advertiserrepository.NewMongoDBAdvertiserRepository(dbStore, logger)
}

func TestSaveAdvertiser(t *testing.T) {
	a := advertiser.NewAdvertiser(shared.NewID(), "Donkey kong drinks")

	err := advertiserRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSavedAdvertiser(t *testing.T) {
	a := advertiser.NewAdvertiser(shared.NewID(), "Donkey kong drinks")

	err := advertiserRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}

	fa,err := advertiserRepo.Get(context.TODO(),a.UserID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fa,a) {
		t.FailNow()
	}
}

func TestGetAllAdveritsers(t *testing.T) {
	a := advertiser.NewAdvertiser(shared.NewID(), "Donkey kong drinks")

	err := advertiserRepo.Save(context.TODO(), a)
	if err != nil {
		t.Fatal(err)
	}

	ads,err := advertiserRepo.Advertisers()
	if err != nil {
		t.Fatal(err)
	}

	if len(ads) == 0 {
		t.FailNow()
	}

}
