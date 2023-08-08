package testdoubles

import (
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
)

type PublisherSignUpUnitOfWorkStub struct{}

func NewPublisherSignUpUnitOfWorkStub() PublisherSignUpUnitOfWorkStub {
	return PublisherSignUpUnitOfWorkStub{}
}

func (p PublisherSignUpUnitOfWorkStub) Save(usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error {
	return nil
}

func (p PublisherSignUpUnitOfWorkStub) UserRepository() user.UserRepository {
	return UserRepositoryStub{}
}
