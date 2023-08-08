package manageduserinteractors

import (
	"errors"

	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type RequestPasswordResetEmailService interface {
	SendMail(u user.User) error
}

type ManagedUserInteractor struct {
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
	tokenGenerator        identityinteractors.AccessTokenGenerator
	emailService          RequestPasswordResetEmailService
}

func NewManagedUserInteractor(
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository,
	tokenGenerator identityinteractors.AccessTokenGenerator,
	emailService RequestPasswordResetEmailService) ManagedUserInteractor {
	return ManagedUserInteractor{
		userRepository,
		managedUserRepository,
		tokenGenerator,
		emailService,
	}
}

var ErrAuthorization = errors.New("authorization error")
var ErrEmailUnverified = errors.New("email is not verified")

func (m *ManagedUserInteractor) SignIn(email user.Email, password string) (string, error) {
	u, err := m.credentialsValid(email, password)
	if err != nil {
		return "", ErrAuthorization
	}

	if !u.IsEmailVerified() {
		return "", ErrEmailUnverified
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

func (m *ManagedUserInteractor) RequestResetPassword(email user.Email) error {

	usr, err := m.userRepository.UserWithEmail(email)
	if err != nil {
		return err
	}

	if !usr.IsEmailVerified() {
		return identityinteractors.ErrEmailNotVerified
	}

	m.emailService.SendMail(usr)
	return nil
}
