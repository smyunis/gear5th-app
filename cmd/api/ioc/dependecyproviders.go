//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-api/cmd/api/controllers/identitycontrollers"
	"gitlab.com/gear5th/gear5th-api/cmd/api/controllers/publishercontrollers"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/mail/identityemail"
)

// API Controllers
func InitManagedUserController() identitycontrollers.ManagedUserController {
	wire.Build(Container)
	return identitycontrollers.ManagedUserController{}
}

func InitPublisherSignUpController() publishercontrollers.PublisherSignUpController {
	wire.Build(Container)
	return publishercontrollers.PublisherSignUpController{}
}

func InitRequestPasswordResetController() identitycontrollers.RequestPasswordResetController {
	wire.Build(Container)
	return identitycontrollers.RequestPasswordResetController{}
}

func InitVerifyEmailController() identitycontrollers.VerifyEmailController {
	wire.Build(Container)
	return identitycontrollers.VerifyEmailController{}
}


// Event Handlers

func InitVerifcationEmailSender() identityemail.VerifcationEmailSender {
	wire.Build(Container)
	return identityemail.VerifcationEmailSender{}
}
