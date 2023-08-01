package shared

import "regexp"

type PhoneNumber struct {
	phoneNumber string
}

func NewPhoneNumber(phoneNumber string) (PhoneNumber, error) {
	regexPattern := `^[\+]?[(]?[0-9]{3}[)]?[-\s\.]?[0-9]{3}[-\s\.]?[0-9]{4,6}$`
	phoneNumberRegex, err := regexp.Compile(regexPattern)
	if err != nil {
		return PhoneNumber{phoneNumber}, InvalidValueError{ValueType: "phone number regex", Value: regexPattern}
	}

	if phoneNumberRegex.MatchString(phoneNumber) {
		return PhoneNumber{phoneNumber}, nil
	}

	return PhoneNumber{phoneNumber}, InvalidValueError{ValueType: "phone number", Value: phoneNumber}
}

func (p PhoneNumber) PhoneNumber() string {
	return p.phoneNumber;
}
