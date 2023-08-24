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
	UserRepository() user.UserRepository
}

type PublisherSignUpInteractor struct {
	eventDispatcher          application.EventDispatcher
	unitOfWork               PublisherSignUpUnitOfWork
	logger                   application.Logger
}

func NewPublisherSignUpInteractor(
	eventDispatcher application.EventDispatcher,
	unitOfWork PublisherSignUpUnitOfWork,
	logger application.Logger) PublisherSignUpInteractor {
	return PublisherSignUpInteractor{
		eventDispatcher,
		unitOfWork,
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
