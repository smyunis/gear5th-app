package publisherinteractors_test

import (
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-api/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	teardown()
}

var interactor publisherinteractors.PublisherSignUpInteractor

func setup() {
	// userRepositoryStub := testdoubles.UserRepositoryStub{}
	// managedUserRepositoryStub := testdoubles.ManagedUserRepositoryStub{}
	// pubRepoStub := testdoubles.PublisherRepositoryStub{}
	pubSignupUnitOfWork := testdoubles.PublisherSignUpUnitOfWorkStub{}
	interactor = publisherinteractors.NewPublisherSignUpInteractor(pubSignupUnitOfWork)
}

func teardown() {
}

func TestPublisherManagedUserSignUp(t *testing.T) {
	email, _ := user.NewEmail("mymail@gmail.com")
	usr := user.NewUser(email)
	manageduser := usr.AsManagedUser(user.NewPersonNameWithFullName("Son Goku"), "gokuisking")

	_ = interactor.ManagedUserSignUp(usr, manageduser)
}





