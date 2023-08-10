//go:build integration
package manageduserrepository_test

import (
	"context"
	"os"
	"reflect"
	"testing"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/identitypersistence/manageduserrepository"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	teardown()
}

var managedUserRepository user.ManagedUserRepository

func setup() {
	configProvider := infrastructure.EnvConfigurationProvider{}
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	managedUserRepository = manageduserrepository.NewMongoDBMangageUserRepository(dbStore)
}

func teardown() {

}

func TestSaveManagedUser(t *testing.T) {
	userEmail, _ := user.NewEmail("mihawk_blade@proton.me")
	u := user.NewUser(userEmail)
	mu := u.AsManagedUser(user.NewPersonNameWithFullName("Dracule Mihawk"), "yoruissharp")

	err := managedUserRepository.Save(context.Background(), mu)
	if err != nil {
		t.Fatal(err.Error())
	}

}

func TestGetWhatWasSaved(t *testing.T) {
	id := shared.ID("id-xxx-yyy")
	m := user.ReconstituteManagedUser(id, user.NewPersonNameWithFullName("Dracule Mihawk"), "pass")

	err := managedUserRepository.Save(context.Background(), m)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmu, err := managedUserRepository.Get(context.Background(), id)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(m, fmu) {
		t.Fatal("fetched not equal to saved")
	}
}
