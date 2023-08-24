package user_test

import (
	"testing"
	"time"

	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
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

	_, ok := events["user/signedup"]

	if !ok {
		t.FailNow()
	}
}

// Assumption is made that this test will not stall for a whole day ;)
func TestAddsSignUpDateForuser(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)

	if u.SignUpDate().Year() != time.Now().Year() {
		t.FailNow()
	}
	if u.SignUpDate().Month() != time.Now().Month() {
		t.FailNow()
	}
	if u.SignUpDate().Day() != time.Now().Day() {
		t.FailNow()
	}
}

func TestReconsUser(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)
	id := u.UserID()

	ru := user.ReconstituteUser(id, userEmail,
		u.PhoneNumber(), u.IsEmailVerified(),
		u.Roles(), u.AuthenticationMethod(), u.SignUpDate())

	if ru.Email() != u.Email() &&
		ru.UserID() != u.UserID() &&
		!ru.SignUpDate().Equal(u.SignUpDate()) {
		t.Fatal("recons not equal")
	}
}

func TestReonsManagedUser(t *testing.T) {
	userEmail, _ := user.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)
	id := u.UserID()
	name := user.NewPersonNameWithFullName("Dracule Mihawk")
	mu := u.AsManagedUser(name, "123456")

	rmu := user.ReconstituteManagedUser(id, name, mu.HashedPassword())

	if rmu.UserID() != mu.UserID() && rmu.Name() != mu.Name() {
		t.Fatal("recons not equal")

	}
}
