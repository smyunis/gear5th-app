package user_test

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"testing"
)

func TestCreateManagedUser(t *testing.T) {
	userEmail, _ := shared.NewEmail("smyunis@outlook.com")
	userPhone, _ := shared.NewPhoneNumber("0929186232")

	u := user.NewUser(userEmail)
	m := u.AsManagedUser("Salman Mohammed", "pass1234", userPhone)

	if u.AuthenticationMethod() != user.Managed {
		t.FailNow()
	}

	if m.AuthenticationMethod() != user.Managed {
		t.FailNow()
	}

}

func TestVerifyManagedUsersEmail(t *testing.T) {
	userEmail, _ := shared.NewEmail("smyunis@outlook.com")
	userPhone, _ := shared.NewPhoneNumber("0929186232")
	u := user.NewUser(userEmail)
	m := u.AsManagedUser("Salman Mohammed", "pass1234", userPhone)

	m.VerifyEmail()

	if !m.IsEmailVerified() {
		t.FailNow()
	}
}

func TestCreateOAuthUser(t *testing.T) {
	userEmail, _ := shared.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)

	o := u.AsOAuthUser("idxxxx-yyyy", user.GoogleOAuth)

	if u.AuthenticationMethod() != user.OAuth {
		t.FailNow()
	}
	if o.AuthenticationMethod() != user.OAuth {
		t.FailNow()
	}
}

func TestSignUpUserAsPublisher(t *testing.T) {
	userEmail, _ := shared.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail)
	o := u.AsOAuthUser("idxxxx-yyyy", user.GoogleOAuth)

	_ = o.SignUpPublisher()

	if !u.HasRole(user.Publisher) {
		t.FailNow()
	}
	if !o.HasRole(user.Publisher) {
		t.FailNow()
	}
}

func TestPasswordForManagedUser(t *testing.T) {
	userEmail, _ := shared.NewEmail("smyunis@outlook.com")
	userPhone, _ := shared.NewPhoneNumber("0929186232")
	u := user.NewUser(userEmail)
	m := u.AsManagedUser("Salman Mohammed", "gokuisking", userPhone)

	m.SetPassword("gokuisking")

	if !m.IsPasswordCorrect("gokuisking") {
		t.FailNow()
	}
}

func TestWrongPasswordForManagedUser(t *testing.T) {
	userEmail, _ := shared.NewEmail("smyunis@outlook.com")
	userPhone, _ := shared.NewPhoneNumber("0929186232")
	u := user.NewUser(userEmail)
	m := u.AsManagedUser("Salman Mohammed", "pass123", userPhone)

	if m.IsPasswordCorrect("gokuisking") {
		t.FailNow()
	}
}

