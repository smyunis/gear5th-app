package publisherinteractors

import (
	"errors"
	"fmt"

	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type PublisherSignUpUnitOfWork interface {
	Save(usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error
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
	existingUser, err := userRepository.UserWithEmail(usr.Email())

	if err == nil {
		usr = existingUser
	} else if !errors.As(err, &shared.ErrEntityNotFound{}) {
		return fmt.Errorf("get user failed: %w", err)
	}
	pub := usr.SignUpPublisher()

	err = i.unitOfWork.Save(usr, managedUser, pub)

	if err != nil {
		return fmt.Errorf("signup publisher failed : %w", err)
	}

	application.ApplicationEventDispatcher.DispatchAsync(usr.DomainEvents())

	return nil

}

// func (i *PublisherSignUpInteractor) savePublisher(usr *user.User, managedUser *user.ManagedUser) error {
// 	pub := usr.SignUpPublisher()

// 	err := i.userRepository.Save(*usr)
// 	if err != nil {
// 		return fmt.Errorf("save user failed: %w", err)
// 	}
// 	err = i.managedUserRepository.Save(*managedUser)
// 	if err != nil {
// 		return fmt.Errorf("save managed user failed : %w", err)
// 	}
// 	err = i.publisherRepository.Save(pub)
// 	if err != nil {
// 		return fmt.Errorf("save publisher failed : %w", err)
// 	}
// 	return nil
// }
