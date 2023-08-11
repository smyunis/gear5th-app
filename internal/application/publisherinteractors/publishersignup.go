package publisherinteractors

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
)

type PublisherSignUpUnitOfWork interface {
	Save(ctx context.Context, usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error
	UserRepository() user.UserRepository
}

type PublisherSignUpInteractor struct {
	unitOfWork PublisherSignUpUnitOfWork
}

func NewPublisherSignUpInteractor(unitOfWork PublisherSignUpUnitOfWork) PublisherSignUpInteractor {
	return PublisherSignUpInteractor{
		unitOfWork,
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

	application.ApplicationEventDispatcher.DispatchAsync(usr.DomainEvents())

	return nil

}
