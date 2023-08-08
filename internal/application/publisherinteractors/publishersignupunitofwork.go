package publisherinteractors

import (


	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
)

type PublisherSignUpUnitOfWork interface {
	Save(usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error
	UserRepository() user.UserRepository
}


// TODO Move this to persistence layer


// type PublisherSignUpUnitOfWorkImpl struct {
// 	userRepository        user.UserRepository
// 	managedUserRepository user.ManagedUserRepository
// 	publisherRepository   publisher.PublisherRepository
// }

// func NewPublisherSignUpUnitOfWorkImpl(
// 	userRepository user.UserRepository,
// 	managedUserRepository user.ManagedUserRepository,
// 	publisherRepository publisher.PublisherRepository) PublisherSignUpUnitOfWorkImpl {
// 	return PublisherSignUpUnitOfWorkImpl{
// 		userRepository,
// 		managedUserRepository,
// 		publisherRepository,
// 	}
// }

// func (uow *PublisherSignUpUnitOfWorkImpl) Save(usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error {
// 	// TODO this method needs to be transactional
// 	err := uow.userRepository.Save(usr)
// 	if err != nil {
// 		return fmt.Errorf("save user failed: %w", err)
// 	}
// 	err = uow.managedUserRepository.Save(managedUser)
// 	if err != nil {
// 		return fmt.Errorf("save managed user failed : %w", err)
// 	}
// 	err = uow.publisherRepository.Save(pub)
// 	if err != nil {
// 		return fmt.Errorf("save publisher failed : %w", err)
// 	}
// 	return nil
// }

// func (uow *PublisherSignUpUnitOfWorkImpl) UserRepository() user.UserRepository {
// 	return uow.userRepository
// }
