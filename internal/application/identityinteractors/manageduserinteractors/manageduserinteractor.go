package manageduserinteractors

import (
	"errors"

	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type ManagedUserInteractor struct {
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
	tokenGenerator        identityinteractors.AccessTokenGenerator
}

func NewManagedUserInteractor(
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository,
	tokenGenerator identityinteractors.AccessTokenGenerator) ManagedUserInteractor {
	return ManagedUserInteractor{
		userRepository:        userRepository,
		managedUserRepository: managedUserRepository,
		tokenGenerator:        tokenGenerator,
	}
}

var ErrAuthorization = errors.New("authorization error")

func (m *ManagedUserInteractor) SignIn(email user.Email, password string) (string, error) {
	u, err := m.credentialsValid(email, password)
	if err != nil {
		return "", ErrAuthorization
	}
	return m.tokenGenerator.Generate(u.UserID())
}



func (m *ManagedUserInteractor) credentialsValid(email user.Email, password string) (user.User, error) {

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
