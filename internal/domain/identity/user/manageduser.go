package user

import "gitlab.com/gear5th/gear5th-api/internal/domain/shared"

type ManagedUser struct {
	User
	fullName        string
	phoneNumber     shared.PhoneNumber
	isEmailVerified bool
}

func (m *ManagedUser) VerifyEmail()  {
	m.isEmailVerified = true
}

func (m *ManagedUser) IsEmailVerified() bool {
	return m.isEmailVerified
}
