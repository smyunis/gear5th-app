package identityusecases_test

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityusecases"
	"gitlab.com/gear5th/gear5th-api/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	teardown()
}

var userRepositoryStub user.UserRepository
var managedUserRepositoryStub user.ManagedUserRepository
var tokenGenerator identityusecases.AccessTokenGenerator

var interactor identityusecases.ManagedUserInteractor

func setup() {
	userRepositoryStub = testdoubles.UserRepositoryStub{}
	managedUserRepositoryStub = testdoubles.ManagedUserRepositoryStub{}
	tokenGenerator = testdoubles.JwtAccessTokenGeneratorStub{}

	interactor = identityusecases.NewManagedUserInteractor(userRepositoryStub, managedUserRepositoryStub, tokenGenerator)
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
