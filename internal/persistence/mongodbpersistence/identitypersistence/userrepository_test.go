//go:build integration

package identitypersistence_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/identitypersistence"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	teardown()
}

var userRepository user.UserRepository

func setup() {
	configProvider := infrastructure.EnvConfigurationProvider{}
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	userRepository = identitypersistence.NewMongoDBUserRepository(dbStore)
}

func teardown() {

}

func TestCanSaveUser(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)

	userRepository.Save(u)
}

func TestCanGetUser(t *testing.T) {
	id := shared.Id("01H7E2AETWXRH2NRZXC47DAK5B")

	u, _ := userRepository.Get(id)

	t.Log(u)
}

func TestCanGetWhatWasSaved(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	ph, _ := user.NewPhoneNumber("0932383239")
	id := shared.NewID()
	bd := time.Date(2023, time.June, 3, 0, 0, 0, 0, time.Local)

	u := user.ReconstituteUser(id, userEmail, ph, true,
		[]user.UserRole{user.Publisher, user.Administrator}, user.Managed, bd)

	err := userRepository.Save(u)
	if err != nil {
		t.Fatal(err.Error())
	}

	fu, err := userRepository.Get(id)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !reflect.DeepEqual(u, fu) {
		t.FailNow()
	}
}
