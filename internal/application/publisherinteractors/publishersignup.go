package publisherinteractors

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
)

type PublisherSignUpUnitOfWork interface {
	Save(ctx context.Context, usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error
	SaveOAuthUser(ctx context.Context, usr user.User, oauthUser user.OAuthUser, pub publisher.Publisher) error
	UserRepository() user.UserRepository
}

type PublisherSignUpInteractor struct {
	eventDispatcher    application.EventDispatcher
	unitOfWork         PublisherSignUpUnitOfWork
	googleOAuthService user.GoogleOAuthService
	logger             application.Logger
}

func NewPublisherSignUpInteractor(
	eventDispatcher application.EventDispatcher,
	unitOfWork PublisherSignUpUnitOfWork,
	googleOAuthService user.GoogleOAuthService,
	logger application.Logger) PublisherSignUpInteractor {
	return PublisherSignUpInteractor{
		eventDispatcher,
		unitOfWork,
		googleOAuthService,
		logger,
	}
}

func (i *PublisherSignUpInteractor) ManagedUserSignUp(usr user.User, managedUser user.ManagedUser) error {

	userRepository := i.unitOfWork.UserRepository()
	existingUser, err := userRepository.UserWithEmail(context.Background(), usr.Email())

	if err == nil {
		usr = existingUser
	} else if !errors.Is(err, application.ErrEntityNotFound) {
		return fmt.Errorf("get user failed: %w", err)
	}

	if usr.HasRole(user.Publisher) {
		return application.ErrConflictFound
	}

	pub := usr.SignUpPublisher()

	err = i.unitOfWork.Save(context.Background(), usr, managedUser, pub)

	if err != nil {
		return fmt.Errorf("signup publisher failed : %w", err)
	}

	i.eventDispatcher.DispatchAsync(usr.DomainEvents())

	return nil

}

func (i *PublisherSignUpInteractor) OAuthUserSignUp(token string) error {
	userDetails, err := i.googleOAuthService.ValidateToken(token)
	if err != nil {
		return err
	}

	userEmail, err := user.NewEmail(userDetails.Email)
	if err != nil {
		return err
	}

	userRepository := i.unitOfWork.UserRepository()
	existingUser, err := userRepository.UserWithEmail(context.Background(), userEmail)
	if err == nil {
		if existingUser.HasRole(user.Publisher) {
			return application.ErrConflictFound
		}
	} else if !errors.Is(err, application.ErrEntityNotFound) {
		return err
	}

	u := user.NewUser(userEmail)
	oauthUser := u.AsOAuthUser(userDetails.AccountID, user.GoogleOAuth)
	p := u.SignUpPublisher()

	err = i.unitOfWork.SaveOAuthUser(context.Background(), u, oauthUser, p)

	if err != nil {
		return err
	}

	i.eventDispatcher.DispatchAsync(u.DomainEvents())

	return nil

}
