package manageduserinteractors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

var ErrInvalidPasswordResetToken = errors.New("reset password token is invalid")

type RequestPasswordResetEmailService interface {
	SendMail(u user.User, resetPasswordToken string) error
}

type ManagedUserInteractor struct {
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
	tokenGenerator        identityinteractors.AccessTokenGenerator
	emailService          RequestPasswordResetEmailService
	kvStore               application.KeyValueStore
}

func NewManagedUserInteractor(
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository,
	tokenGenerator identityinteractors.AccessTokenGenerator,
	emailService RequestPasswordResetEmailService,
	kvStore application.KeyValueStore) ManagedUserInteractor {
	return ManagedUserInteractor{
		userRepository,
		managedUserRepository,
		tokenGenerator,
		emailService,
		kvStore,
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

	u, err := m.userRepository.UserWithEmail(context.Background(), email)
	if err != nil {
		return u, ErrAuthorization
	}

	managedUser, err := m.managedUserRepository.Get(context.Background(), u.UserID())
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

	usr, err := m.userRepository.UserWithEmail(context.Background(), email)
	if err != nil {
		return err
	}

	if !usr.IsEmailVerified() {
		return identityinteractors.ErrEmailNotVerified
	}

	//Using shared.NewID()  generators abitlity to generate random strings to be used as token
	token := shared.NewID().String()
	err = m.kvStore.Save(m.passwordResetStorageKey(usr), token, 30*time.Minute)
	if err != nil {
		return err
	}

	m.emailService.SendMail(usr, token)
	return nil
}

func (m *ManagedUserInteractor) ResetPassword(email user.Email, newPassword, resetToken string) error {

	u, err := m.userRepository.UserWithEmail(context.Background(), email)
	if err != nil {
		return err
	}

	if !u.IsEmailVerified() {
		return identityinteractors.ErrEmailNotVerified
	}

	token, err := m.kvStore.Get(m.passwordResetStorageKey(u))
	if err != nil {
		return fmt.Errorf("reset password failed: %w", err)
	}

	if token != resetToken {
		return ErrInvalidPasswordResetToken
	}

	managedUser, err := m.managedUserRepository.Get(context.Background(), u.UserID())
	if err != nil {
		return err
	}

	err = managedUser.SetPassword(newPassword)
	if err != nil {
		return err
	}

	return nil
}

func (ManagedUserInteractor) passwordResetStorageKey(u user.User) string {
	return fmt.Sprintf("identity:manageduser:%s:passwordresettoken", u.UserID().String())
}
