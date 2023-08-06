package user_test

import (
	"testing"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

func TestInvalidPhonenumber(t *testing.T) {
	invalidPhoneNumber := "hsdjfh"

	_, err := user.NewPhoneNumber(invalidPhoneNumber)

	if err == nil {
		t.Fatal(err.Error())
	}
}

func TestValidPhoneNumber(t *testing.T) {
	validPhones := []string{"0929186232", "+251929186232", "0799116232",
		"+251799116232", "251799116232", "251929186232"}

	for _, phonenum := range validPhones {
		t.Run(phonenum, func(t *testing.T) {
			phnum, err := user.NewPhoneNumber(phonenum)
			if err != nil {
				t.Fatal(err.Error())
			}
			if phnum.PhoneNumber() != phonenum {
				t.FailNow()
			}
		})
	}
}
