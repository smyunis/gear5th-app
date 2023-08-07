//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-api/cmd/api/identitycontrollers"
	"gitlab.com/gear5th/gear5th-api/cmd/api/publishercontrollers"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/publisher/publisher"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/identity/accesstoken"
)

var Container wire.ProviderSet = wire.NewSet(
	//Stub Repositories
	wire.Struct(new(testdoubles.UserRepositoryStub), "*"),
	wire.Struct(new(testdoubles.ManagedUserRepositoryStub), "*"),
	wire.Struct(new(testdoubles.PublisherRepositoryStub), "*"),

	// Repositories
	wire.Bind(new(user.UserRepository), new(testdoubles.UserRepositoryStub)),
	wire.Bind(new(user.ManagedUserRepository), new(testdoubles.ManagedUserRepositoryStub)),
	wire.Bind(new(publisher.PublisherRepository), new(testdoubles.PublisherRepositoryStub)),

	//Infrastructures
	accesstoken.NewJwtAccessTokenGenenrator,
	wire.Bind(new(identityinteractors.AccessTokenGenerator), new(accesstoken.JwtAccessTokenGenenrator)),

	//Interactors
	manageduserinteractors.NewManagedUserInteractor,
	publisherinteractors.NewPublisherSignUpInteractor,

	//Controllers
	publishercontrollers.NewPublisherSignUpController,
	identitycontrollers.NewManagedUserController)
