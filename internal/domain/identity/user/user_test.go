package user_test

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"testing"
)

func TestCreateManagedUser(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")

	u := user.NewUser(userEmail)
	u.AsManagedUser(user.NewPersonNameWithFullName("Salman Mohammed"), "pass1234")

	if u.AuthenticationMethod() != user.Managed {
		t.FailNow()
	}

}

func TestVerifyManagedUsersEmail(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)

	u.VerifyEmail()

	if !u.IsEmailVerified() {
		t.FailNow()
	}
}

func TestCreateOAuthUser(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)

	u.AsOAuthUser("idxxxx-yyyy", user.GoogleOAuth)

	if u.AuthenticationMethod() != user.OAuth {
		t.FailNow()
	}

}

func TestSignUpUserAsPublisher(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)
	_ = u.SignUpPublisher()

	if !u.HasRole(user.Publisher) {
		t.FailNow()
	}

}

func TestPasswordForManagedUser(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)
	m := u.AsManagedUser(user.NewPersonNameWithFullName("Salman Mohammed"), "gokuisking")

	m.SetPassword("vegetaisking")

	if !m.IsPasswordCorrect("vegetaisking") {
		t.FailNow()
	}
}

func TestWrongPasswordForManagedUser(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)
	m := u.AsManagedUser(user.NewPersonNameWithFullName("Salman Mohammed"), "pass123")

	if m.IsPasswordCorrect("gokuisking") {
		t.FailNow()
	}
}

func TestCreateNewUserEmitsDomainEvent(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)

	events := u.DomainEvents()

	_, ok := events["user.created"]

	if !ok {
		t.FailNow()
	}
}
