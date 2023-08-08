package user

import (
	"regexp"

	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	regexPattern := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,6}$`
	emailRegex, err := regexp.Compile(regexPattern)
	if err != nil {
		return Email{email}, shared.ErrInvalidValue{ValueType: "email regex", Value: regexPattern}
	}

	if emailRegex.MatchString(email) {
		return Email{email: email}, nil
	}

	return Email{email}, shared.ErrInvalidValue{ValueType: "email", Value: email}

}

func (e Email) Email() string {
	return e.email
}
