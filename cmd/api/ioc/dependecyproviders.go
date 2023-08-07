//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-api/cmd/api/identitycontrollers"
	"gitlab.com/gear5th/gear5th-api/cmd/api/publishercontrollers"
)

func InitManagedUserController() identitycontrollers.ManagedUserController {
	wire.Build(Container)
	return identitycontrollers.ManagedUserController{}
}

func InitPublisherSignUpController() publishercontrollers.PublisherSignUpController {
	wire.Build(Container)
	return publishercontrollers.PublisherSignUpController{}
}
