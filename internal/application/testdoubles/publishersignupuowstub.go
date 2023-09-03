package testdoubles

import (
	"context"

	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/publisher"
)

type PublisherSignUpUnitOfWorkStub struct{}

// SaveOAuthUser implements publisherinteractors.PublisherSignUpUnitOfWork.
func (PublisherSignUpUnitOfWorkStub) SaveOAuthUser(ctx context.Context, usr user.User, oauthUser user.OAuthUser, pub publisher.Publisher) error {
	return nil
}

func NewPublisherSignUpUnitOfWorkStub() PublisherSignUpUnitOfWorkStub {
	return PublisherSignUpUnitOfWorkStub{}
}


func (p PublisherSignUpUnitOfWorkStub) Save(ctx context.Context, usr user.User, managedUser user.ManagedUser, pub publisher.Publisher) error {
	return nil
}

func (p PublisherSignUpUnitOfWorkStub) UserRepository() user.UserRepository {
	return UserRepositoryStub{}
}
