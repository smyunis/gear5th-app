package application

import "gitlab.com/gear5th/gear5th-app/internal/domain/shared"

type AccessTokenService interface {
	Generate(subject shared.ID) (string, error)
	Validate(token string) bool
	UserID(token string) (shared.ID, error)
}

// type AccessTokenService interface {

// }

type DigitalSignatureService interface {
	Generate(message string) (string, error)
	Validate(hmacMessage string) bool
	GetMessage(hashed string) (string, error)
}
