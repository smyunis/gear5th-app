package testdoubles

import (
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
)

type DigitalSignatureValidationServiceMock struct {
	msg string
}

var dsvMsg string

func (m *DigitalSignatureValidationServiceMock) Generate(message string) (string, error) {
	m.msg = message
	return message + " xxx", nil
}

func (m *DigitalSignatureValidationServiceMock) Validate(hmacMessage string) bool {
	return true
}

func (m *DigitalSignatureValidationServiceMock) GetMessage(message string) (string, error) {
	return strings.Split(message, " ")[0], nil
}

type GoogleOAuthServiceStub struct{}

func (GoogleOAuthServiceStub) ValidateToken(identityToken string) (user.GoogleOAuthUserDetails, error) {
	return user.GoogleOAuthUserDetails{
		"Salman Yunis",
		StubID,
		"doni793doni793@gmail.com",
	}, nil
}
