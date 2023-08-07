package identityinteractors

import "gitlab.com/gear5th/gear5th-api/internal/domain/shared"

type AccessTokenGenerator interface {
	Generate(subject shared.Id) (string, error)
}

type AccessTokenValidator interface {
	Validate(token string) error
}
