package testdoubles

import (
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
)

type VerificationEmailServiceMock struct {
}

func (m VerificationEmailServiceMock) SendMail(u user.User) error {
	return nil
}
