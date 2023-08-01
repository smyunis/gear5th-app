package shared

import "regexp"

type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	regexPattern := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	emailRegex, err := regexp.Compile(regexPattern)
	if err != nil {
		return Email{email}, InvalidValueError{ValueType: "email regex", Value: regexPattern}
	}

	if emailRegex.MatchString(email) {
		return Email{email: email}, nil
	}

	return Email{email}, InvalidValueError{ValueType: "email", Value: email}

}

func (e Email) Email() string {
	return e.email
}
