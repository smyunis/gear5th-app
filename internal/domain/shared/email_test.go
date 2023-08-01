package shared_test

import (
	"testing"

	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

func TestInvalid_email(t *testing.T) {
	invalidEmail := `23doni793doni793gmail.comnet`

	_, err := shared.NewEmail(invalidEmail)

	if err == nil {
		t.Fatal(err.Error())
	}
}

func TestValidEmail(t *testing.T) {
	validEmail := `doni793doni793@gmail.com`

	mail, err := shared.NewEmail(validEmail)

	if err != nil {
		t.Fatal(err.Error())
	}
	if mail.Email() != validEmail {
		t.Fatal("email not saved")
	}
}
