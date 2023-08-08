package user

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"golang.org/x/crypto/bcrypt"
)

type ManagedUserRepository interface {
	shared.EntityRepository[ManagedUser]
}

type ManagedUser struct {
	userId         shared.Id
	name           PersonName
	hashedPassword string
}

func (m *ManagedUser) SetPassword(plainPassword string) error {
	hashed, err := m.hashPlainPassword(plainPassword)
	if err != nil {
		return shared.ErrInvalidValue{ValueType: "password", Value: plainPassword, InnerError: err}
	}
	m.hashedPassword = string(hashed)
	return nil
}

func (m *ManagedUser) IsPasswordCorrect(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.hashedPassword), []byte(plainPassword))
	return err == nil
}

func (m *ManagedUser) hashPlainPassword(plainPassword string) ([]byte, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashBytes, nil
}
