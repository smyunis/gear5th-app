package testdoubles

import "strings"

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
	return strings.Split(message, " ")[0],nil
}
