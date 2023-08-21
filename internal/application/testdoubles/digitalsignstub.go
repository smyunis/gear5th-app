package testdoubles

type DigitalSignatureValidationServiceMock struct{}

func (m DigitalSignatureValidationServiceMock) Generate(message string) (string, error) {
	return "token", nil
}

func (m DigitalSignatureValidationServiceMock) Validate(hmacMessage string) bool {
	return true
}

func (m DigitalSignatureValidationServiceMock) GetMessage(message string) (string, error) {
	return "msg",nil
}