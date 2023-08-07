//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-api/cmd/api/identitycontrollers"
)



func InitManagedUserController() identitycontrollers.ManagedUserController {
	wire.Build(Container)
	return identitycontrollers.ManagedUserController{}
}
