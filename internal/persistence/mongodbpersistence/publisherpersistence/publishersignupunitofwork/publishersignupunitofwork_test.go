//go:build integration
package publishersignupunitofwork_test

import (
	"context"
	"os"
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/identitypersistence/manageduserrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/identitypersistence/userrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/publisherpersistence/publisherrepository"
	"gitlab.com/gear5th/gear5th-api/internal/persistence/mongodbpersistence/publisherpersistence/publishersignupunitofwork"
)

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())

	teardown()
}

var pubSignUpUnitOfWork publisherinteractors.PublisherSignUpUnitOfWork

func setup() {
	configProvider := infrastructure.EnvConfigurationProvider{}
	dbStore := mongodbpersistence.NewMongoDBStoreBootstrap(configProvider)
	userRepository := userrepository.NewMongoDBUserRepository(dbStore)
	managedUserRepository := manageduserrepository.NewMongoDBMangageUserRepository(dbStore)
	publisherRepository := publisherrepository.NewMongoDBPublisherRepository(dbStore)

	pubSignUpUnitOfWork = publishersignupunitofwork.NewMongoDBPublisherSignUpUnitOfWork(dbStore, userRepository, managedUserRepository, publisherRepository)

}

func teardown() {

}

func TestSave(t *testing.T) {
	userEmail, _ := user.NewEmail("redhair@outlook.com")
	ph, _ := user.NewPhoneNumber("0977231546")
	id := shared.ID("idx-xxx-yyy")
	bd := time.Date(2025, time.June, 3, 0, 0, 0, 0, time.Local)
	u := user.ReconstituteUser(id, userEmail, ph, true,
		[]user.UserRole{user.Publisher, user.Administrator}, user.Managed, bd)

	m := user.ReconstituteManagedUser(id, user.NewPersonNameWithFullName("Akagami Shanks"), "pass")

	notifications := []shared.Notification{
		shared.ReconstituteNotification("N1", time.Date(2023, time.November, 5, 0, 0, 0, 0, time.Local)),
		shared.ReconstituteNotification("N2", time.Date(2022, time.July, 5, 0, 0, 0, 0, time.Local)),
	}
	pub := publisher.ReconstitutePublisher(id, notifications)

	err := pubSignUpUnitOfWork.Save(context.Background(), u, m, pub)

	if err != nil {
		t.Fatal(err.Error())
	}

}
