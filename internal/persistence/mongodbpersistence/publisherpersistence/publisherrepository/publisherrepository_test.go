//go:build integration
package publisherrepository_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/publisherpersistence/publisherrepository"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	teardown()
}

var publisherRepository publisher.PublisherRepository

func setup() {
	configProvider := infrastructure.EnvConfigurationProvider{}
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	publisherRepository = publisherrepository.NewMongoDBPublisherRepository(dbStore)
}

func teardown() {

}

func TestCanSavePublisher(t *testing.T) {
	pub := publisher.NewPublisher(shared.NewID())
	pub.Notify(shared.NewNotification("Alert! Some stuff needs your attention"))

	err := publisherRepository.Save(context.Background(), pub)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCanGetSavedPublisher(t *testing.T) {
	notifications := []shared.Notification{
		shared.ReconstituteNotification("N1", time.Date(2023, time.November, 5, 0, 0, 0, 0, time.Local)),
		shared.ReconstituteNotification("N2", time.Date(2022, time.July, 5, 0, 0, 0, 0, time.Local)),
	}
	pub := publisher.ReconstitutePublisher(shared.NewID(), notifications)

	err := publisherRepository.Save(context.Background(), pub)
	if err != nil {
		t.Fatal(err.Error())
	}

	fp, err := publisherRepository.Get(context.Background(), pub.UserID())
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(fp, pub) {
		t.Log(fp)
		t.Log(pub)
		t.FailNow()
	}
}
