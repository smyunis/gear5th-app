package publisherinteractors

import (
	"fmt"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
)

type PublisherSignUpInteractor struct {
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
	publisherRepository   publisher.PublisherRepository
}

func NewPublisherSignUpInteractor(
	userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository, publisherRepository publisher.PublisherRepository) PublisherSignUpInteractor {
	return PublisherSignUpInteractor{
		userRepository,
		managedUserRepository,
		publisherRepository,
	}
}

func (i *PublisherSignUpInteractor) ManagedUserSignUp(usr *user.User, managedUser *user.ManagedUser) error {

	pub := usr.SignUpPublisher()

	err := i.userRepository.Save(*usr)
	if err != nil {
		return fmt.Errorf("save user failed: %w", err)
	}
	err = i.managedUserRepository.Save(*managedUser)
	if err != nil {
		return fmt.Errorf("save managed user failed : %w", err)
	}
	err = i.publisherRepository.Save(pub)
	if err != nil {
		return fmt.Errorf("save publisher failed : %w", err)
	}

	return nil
}
