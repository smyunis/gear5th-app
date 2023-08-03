package user_test

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"testing"
)

func TestCreateManagedUser(t *testing.T) {
	userEmail, _ := shared.NewEmail("smyunis@outlook.com")
	userPhone, _ := shared.NewPhoneNumber("0929186232")

	u := user.NewUser(userEmail, user.Managed)
	m := u.AsManagedUser("Salman Mohammed", userPhone)

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
	u := user.NewUser(userEmail, user.Managed)
	m := u.AsManagedUser("Salman Mohammed", userPhone)

	m.VerifyEmail()

	if !m.IsEmailVerified() {
		t.FailNow()
	}
}

func TestCreateOAuthUser(t *testing.T) {
	userEmail, _ := shared.NewEmail("smyunis@outlook.com")
	u := user.NewUser(userEmail, user.Managed)

	o := u.AsOAuthUser("idxxxx-yyyy", user.GoogleOAuth)

	if u.AuthenticationMethod() != user.OAuth {
		t.FailNow()
	}

	if o.AuthenticationMethod() != user.OAuth {
		t.FailNow()
	}
}
