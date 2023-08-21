package manageduserinteractors

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type RequestPasswordResetEmailService interface {
	SendMail(u user.User, resetPasswordToken string) error
}

type ManagedUserInteractor struct {
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
	tokenGenerator        identityinteractors.AccessTokenGenerator
	emailService          RequestPasswordResetEmailService
	signService           identityinteractors.DigitalSignatureValidationService
}

func NewManagedUserInteractor(
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository,
	tokenGenerator identityinteractors.AccessTokenGenerator,
	emailService RequestPasswordResetEmailService,
	signService identityinteractors.DigitalSignatureValidationService) ManagedUserInteractor {
	return ManagedUserInteractor{
		userRepository,
		managedUserRepository,
		tokenGenerator,
		emailService,
		signService,
	}
}

var ErrAuthorization = errors.New("authorization error")

func (m *ManagedUserInteractor) SignIn(email user.Email, password string) (string, error) {
	u, err := m.CredentialsValid(email, password)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			return "", ErrAuthorization
		}
		return "", err
	}

	if !u.IsEmailVerified() {
		return "", identityinteractors.ErrEmailNotVerified
	}

	return m.tokenGenerator.Generate(u.UserID())
}

func (m *ManagedUserInteractor) CredentialsValid(email user.Email, password string) (user.User, error) {

	u, err := m.userRepository.UserWithEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			return u, ErrAuthorization
		}
		return u, err
	}

	managedUser, err := m.managedUserRepository.Get(context.Background(), u.UserID())
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			return u, ErrAuthorization
		}
		return u, err
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
	token, err = m.signService.Generate(token)
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

	if !m.signService.Validate(resetToken) {
		return ErrInvalidToken
	}

	managedUser, err := m.managedUserRepository.Get(context.Background(), u.UserID())
	if err != nil {
		return err
	}

	err = managedUser.SetPassword(newPassword)
	if err != nil {
		return err
	}

	err = m.managedUserRepository.Save(context.Background(), managedUser)
	if err != nil {
		return err
	}

	return nil
}

func (m *ManagedUserInteractor) VerifyEmail(userId shared.ID, token string) error {

	if !m.signService.Validate(token) {
		return ErrInvalidToken
	}

	u, err := m.userRepository.Get(context.Background(), userId)
	if err != nil {
		return application.ErrEntityNotFound
	}

	u.VerifyEmail()

	err = m.userRepository.Save(context.Background(), u)
	if err != nil {
		return err
	}

	return nil
}

func PasswordResetTokenStoreKey(userId string) string {
	return fmt.Sprintf("identity:user:%s:passwordresettoken", userId)
}

func EmailVerificationTokenStoreKey(userId string) string {
	return fmt.Sprintf("identity:user:%s:emailverificationtoken", userId)
}
