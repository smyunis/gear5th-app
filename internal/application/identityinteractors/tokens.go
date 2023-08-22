package identityinteractors

import "gitlab.com/gear5th/gear5th-app/internal/domain/shared"

type AccessTokenGenerator interface {
	Generate(subject shared.ID) (string, error)
}

type AccessTokenValidator interface {
	Validate(token string) error
}

type DigitalSignatureValidationService interface {
	Generate(message string) (string, error)
	Validate(hmacMessage string) bool
	GetMessage(hashed string) (string, error)
}
