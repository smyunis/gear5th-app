package publisherinteractors_test

import (
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
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
	verfEmailServiceMock := testdoubles.VerificationEmailServiceMock{}
	consoleLogger := testdoubles.ConsoleLogger{}
	interactor = publisherinteractors.NewPublisherSignUpInteractor(pubSignupUnitOfWork, verfEmailServiceMock, consoleLogger)
}

func teardown() {
}

func TestPublisherManagedUserSignUp(t *testing.T) {
	email, _ := user.NewEmail("mymail@gmail.com")
	usr := user.NewUser(email)
	manageduser := usr.AsManagedUser(user.NewPersonNameWithFullName("Son Goku"), "gokuisking")

	_ = interactor.ManagedUserSignUp(usr, manageduser)
}
