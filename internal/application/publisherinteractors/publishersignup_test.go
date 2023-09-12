package publisherinteractors_test

import (
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
)

var publisherInteractor publisherinteractors.PublisherSignUpInteractor

func publisherTestSetup() {
	var evtDispather application.EventDispatcher = &testdoubles.LocalizedEventDispatcher{}
	pubSignupUnitOfWork := testdoubles.PublisherSignUpUnitOfWorkStub{}
	consoleLogger := testdoubles.ConsoleLogger{}
	gStub := testdoubles.GoogleOAuthServiceStub{}
	publisherInteractor = publisherinteractors.NewPublisherSignUpInteractor(evtDispather, pubSignupUnitOfWork, gStub, consoleLogger)
}

func TestPublisherManagedUserSignUp(t *testing.T) {
	email, _ := user.NewEmail("mymail@gmail.com")
	usr := user.NewUser(email)
	manageduser := usr.AsManagedUser(user.NewPersonNameWithFullName("Son Goku"), "gokuisking")

	_ = publisherInteractor.ManagedUserSignUp(usr, manageduser)
}
