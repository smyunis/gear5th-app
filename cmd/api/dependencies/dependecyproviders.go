//go:build wireinject

package dependencies

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-api/cmd/api/identitycontrollers"
	"gitlab.com/gear5th/gear5th-api/internal/application/identity/usersignin"
	"gitlab.com/gear5th/gear5th-api/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/identity/accesstoken"
)

var Container wire.ProviderSet = wire.NewSet(
	wire.Struct(new(testdoubles.UserRepositoryStub), "*"),
	wire.Struct(new(testdoubles.ManagedUserRepositoryStub), "*"),

	wire.Bind(new(user.UserRepository), new(testdoubles.UserRepositoryStub)),
	wire.Bind(new(user.ManagedUserRepository), new(testdoubles.ManagedUserRepositoryStub)),

	accesstoken.NewJwtAccessTokenGenenrator,
	wire.Bind(new(usersignin.AccessTokenGenerator), new(accesstoken.JwtAccessTokenGenenrator)),
	usersignin.NewManagedUserInteractor,
	identitycontrollers.NewManagedUserSignIn)

func InitManagedUserSignInController() identitycontrollers.ManagedUserSignInController {
	wire.Build(Container)
	return identitycontrollers.ManagedUserSignInController{}
}
