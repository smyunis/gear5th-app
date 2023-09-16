package identityinteractors

import (
	"context"
	"errors"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type RequestPasswordResetEmailService interface {
	SendMail(u user.User, resetPasswordToken string) error
}

type ManagedUserInteractor struct {
	eventDispatcher       application.EventDispatcher
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
	tokenGenerator        application.AccessTokenService
	emailService          RequestPasswordResetEmailService
	signService           application.DigitalSignatureService
}

func NewManagedUserInteractor(
	eventDispatcher application.EventDispatcher,
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository,
	tokenGenerator application.AccessTokenService,
	emailService RequestPasswordResetEmailService,
	signService application.DigitalSignatureService) ManagedUserInteractor {
	return ManagedUserInteractor{
		eventDispatcher,
		userRepository,
		managedUserRepository,
		tokenGenerator,
		emailService,
		signService,
	}
}


func (m *ManagedUserInteractor) SignIn(email user.Email, password string) (string, error) {
	u, err := m.CredentialsValid(email, password)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			return "", application.ErrAuthorization
		}
		return "", err
	}

	if !u.IsEmailVerified {
		return "", ErrEmailNotVerified
	}

	return m.tokenGenerator.Generate(u.ID)
}

func (m *ManagedUserInteractor) CredentialsValid(email user.Email, password string) (user.User, error) {

	u, err := m.userRepository.UserWithEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			return u, application.ErrAuthorization
		}
		return u, err
	}

	managedUser, err := m.managedUserRepository.Get(context.Background(), u.ID)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			return u, application.ErrAuthorization
		}
		return u, err
	}

	if email != u.Email {
		return u, application.ErrAuthorization
	}

	if !managedUser.IsPasswordCorrect(password) {
		return u, application.ErrAuthorization
	}

	return u, nil
}

func (m *ManagedUserInteractor) RequestResetPassword(email user.Email) error {

	usr, err := m.userRepository.UserWithEmail(context.Background(), email)
	if err != nil {
		return err
	}

	if !usr.IsEmailVerified {
		return ErrEmailNotVerified
	}

	token, err := m.signService.Generate(email.String())
	if err != nil {
		return err
	}

	err = m.emailService.SendMail(usr, token)
	if err != nil {
		return err
	}
	return nil
}

func (m *ManagedUserInteractor) ResetPassword(email user.Email, newPassword, resetToken string) error {

	hashedEmail, err := m.signService.GetMessage(resetToken)
	if hashedEmail != email.String() || err != nil {
		return ErrInvalidToken
	}

	if !m.signService.Validate(resetToken) {
		return ErrInvalidToken
	}

	u, err := m.userRepository.UserWithEmail(context.Background(), email)
	if err != nil {
		return err
	}

	if !u.IsEmailVerified {
		return ErrEmailNotVerified
	}

	managedUser, err := m.managedUserRepository.Get(context.Background(), u.ID)
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

func (m *ManagedUserInteractor) VerifyEmail(userID shared.ID, token string) error {

	tokenUserID, err := m.signService.GetMessage(token)
	if err != nil {
		return ErrInvalidToken
	}

	if tokenUserID != userID.String() {
		return ErrInvalidToken
	}

	if !m.signService.Validate(token) {
		return ErrInvalidToken
	}

	u, err := m.userRepository.Get(context.Background(), userID)
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
