package identityinteractors_test

import (
	"errors"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	teardown()
}

var userRepositoryStub user.UserRepository
var managedUserRepositoryStub user.ManagedUserRepository
var tokenGenerator application.AccessTokenService
var kvstore = testdoubles.KVStoreMock{}
var digiSignService = &testdoubles.DigitalSignatureValidationServiceMock{}
var evtDispather application.EventDispatcher = &testdoubles.LocalizedEventDispatcher{}

var interactor identityinteractors.ManagedUserInteractor

func setup() {
	userRepositoryStub = testdoubles.UserRepositoryStub{}
	managedUserRepositoryStub = testdoubles.ManagedUserRepositoryStub{}
	tokenGenerator = testdoubles.JwtAccessTokenGeneratorStub{}
	emailServiceStub := testdoubles.RequestResetPasswordEmailStub{}

	interactor = identityinteractors.NewManagedUserInteractor(
		evtDispather,
		userRepositoryStub,
		managedUserRepositoryStub,
		tokenGenerator, emailServiceStub, digiSignService)
}

func teardown() {
	userRepositoryStub = nil
}

func TestMangedUserSignIn(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")

	interactor.SignIn(mymail, "gokuisking")
}

func TestManagedUserSignInReturnsTokenForValidMail(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")

	token, _ := interactor.SignIn(mymail, "gokuisking")

	if token == "" {
		t.FailNow()
	}
}

func TestManagedUserSignInReturnsErrorForInValidMail(t *testing.T) {
	mymail, _ := user.NewEmail("yourmail@gmail.com")

	_, err := interactor.SignIn(mymail, "gokuisking")

	if err == nil {
		t.FailNow()
	}
}

func TestManagedUserSignInReturnsTokenNoErrForValidEmailPassword(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")

	token, err := interactor.SignIn(mymail, "gokuisking")

	if token == "" {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}
}

func TestManagedUserSignInReturnsErrForValidEmailInvalidPassword(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")

	_, err := interactor.SignIn(mymail, "vegetaisking")

	if err == nil {
		t.FailNow()
	}
}

func TestManagedUserSignInReturnsValidToken(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")
	token, err := interactor.SignIn(mymail, "gokuisking")
	if err != nil {
		t.Fatal(err.Error())
	}
	parser := jwt.NewParser()

	tok, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		t.Fatal(err.Error())
	}

	if iss, err := tok.Claims.GetIssuer(); err != nil || iss != "api.gear5th.com" {
		t.FailNow()
	}
	if sub, err := tok.Claims.GetSubject(); err != nil || sub != "stub-id-xxx" {
		t.FailNow()
	}
	if aud, err := tok.Claims.GetAudience(); err != nil || aud[0] != "api.gear5th.com" {
		t.FailNow()
	}
}

func TestOnlyUsersWithVerifiedEmailsCanSignIn(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")
	_, err := interactor.SignIn(mymail, "gokuisking")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestResetPasswordRequest(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")
	_ = interactor.RequestResetPassword(mymail)
}

func TestResetPasswordRequestForExistingEmail(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")
	err := interactor.RequestResetPassword(mymail)

	if err != nil {
		t.FailNow()
	}
}

func TestResetPasswordRequestForNonExistingEmail(t *testing.T) {
	mymail, _ := user.NewEmail("yourmail@gmail.com")
	err := interactor.RequestResetPassword(mymail)

	if err == nil {
		t.FailNow()
	}
}

func TestResetPasswordRequestFailsForUnverifiedEmail(t *testing.T) {
	mymail, _ := user.NewEmail("somemail@gmail.com")
	err := interactor.RequestResetPassword(mymail)

	if !errors.Is(err, identityinteractors.ErrEmailNotVerified) {
		t.FailNow()
	}
}

func TestResetPasswordRequestFailsForInvalidEmail(t *testing.T) {
	mymail, _ := user.NewEmail("invalidemailaddress.com")
	err := interactor.RequestResetPassword(mymail)

	if err == nil {
		t.FailNow()
	}
}

func TestResetPasswordRequestEmailIsSent(t *testing.T) {

	var emailServiceSpy = testdoubles.RequestResetPasswordEmailSpy{}
	testdoubles.RequestResetPasswordEmailSpyReset()

	interactor := identityinteractors.NewManagedUserInteractor(
		evtDispather,
		userRepositoryStub,
		managedUserRepositoryStub,
		tokenGenerator, emailServiceSpy, digiSignService)

	mymail, _ := user.NewEmail("mymail@gmail.com")

	err := interactor.RequestResetPassword(mymail)

	if err != nil {
		t.FailNow()
	}
	if _, called := testdoubles.RequestResetPasswordEmailSpyGet(); !called {
		t.FailNow()
	}
}

func TestResetPasswordRequestEmailIsNotSentForUnknownEmail(t *testing.T) {

	var emailServiceSpy = testdoubles.RequestResetPasswordEmailSpy{}
	testdoubles.RequestResetPasswordEmailSpyReset()

	digiSignService := &testdoubles.DigitalSignatureValidationServiceMock{}

	interactor := identityinteractors.NewManagedUserInteractor(
		evtDispather,
		userRepositoryStub,
		managedUserRepositoryStub,
		tokenGenerator, emailServiceSpy, digiSignService)

	mymail, _ := user.NewEmail("yourmail@gmail.com")

	err := interactor.RequestResetPassword(mymail)

	if err == nil {
		t.FailNow()
	}
	if _, called := testdoubles.RequestResetPasswordEmailSpyGet(); called {
		t.FailNow()
	}
}

func TestResetPassword(t *testing.T) {
	resetToken := "mypasswordresettoken"
	kvstore.Save("identity:manageduser:passwordresettoken", resetToken, 0)
	mymail, _ := user.NewEmail("mymail@gmail.com")

	interactor.ResetPassword(mymail, "newpass", resetToken)
}

func TestResetPasswordFailsForUnknownEmail(t *testing.T) {
	resetToken := "yourmail@gmail.com xxx"
	mymail, _ := user.NewEmail("yourmail@gmail.com")

	err := interactor.ResetPassword(mymail, "newpass", resetToken)

	if !errors.Is(err, application.ErrEntityNotFound) {
		t.Fatal(err.Error())
	}
}

func TestResetPasswordFailsForUnverifiedEmail(t *testing.T) {
	resetToken := "somemail@gmail.com xxx"
	mymail, _ := user.NewEmail("somemail@gmail.com")

	err := interactor.ResetPassword(mymail, "newpass", resetToken)

	if !errors.Is(err, identityinteractors.ErrEmailNotVerified) {
		t.Fatal(err.Error())
	}
}

func TestResetPasswordFailsForUnknownResetToken(t *testing.T) {
	mymail, _ := user.NewEmail("mymail@gmail.com")
	err := interactor.ResetPassword(mymail, "newpass", "yourpasswordresettoken")
	if err == nil {
		t.Fatal(err.Error())
	}
}

func TestResetPasswordForValidToken(t *testing.T) {
	resetToken := "mymail@gmail.com xxx"
	mymail, _ := user.NewEmail("mymail@gmail.com")
	err := interactor.ResetPassword(mymail, "newpass", resetToken)
	if err != nil {
		t.Fatal(err.Error())
	}
}
