package usersignin

import (
	"errors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type AccessTokenGenerator interface {
	Generate(subject shared.Id) (string, error)
}

type ManagedUserInteractor struct {
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
	tokenGenerator        AccessTokenGenerator
}

func NewManagedUserInteractor(
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository,
	tokenGenerator AccessTokenGenerator) ManagedUserInteractor {
	return ManagedUserInteractor{
		userRepository:        userRepository,
		managedUserRepository: managedUserRepository,
		tokenGenerator:        tokenGenerator,
	}
}

var ErrAuthorization = errors.New("authorization error")

func (m *ManagedUserInteractor) SignIn(email user.Email, password string) (string, error) {

	u, err := m.CredentialsValid(email, password)

	if err != nil {
		return "", ErrAuthorization
	}

	return m.tokenGenerator.Generate(u.UserID())
}

func (m *ManagedUserInteractor) CredentialsValid(email user.Email, password string) (user.User, error) {

	u, err := m.userRepository.UserWithEmail(email)
	if err != nil {
		return u, ErrAuthorization
	}

	managedUser, err := m.managedUserRepository.Get(u.UserID())
	if err != nil {
		return u, ErrAuthorization
	}

	if email != u.Email() {
		return u, ErrAuthorization
	}

	if !managedUser.IsPasswordCorrect(password) {
		return u, ErrAuthorization
	}

	return u, nil
}
