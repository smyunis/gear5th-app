//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-api/web/controllers/publish/identitycontrollers"
	"gitlab.com/gear5th/gear5th-api/web/controllers/publish/publishercontrollers"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure/mail/identityemail"
)

// API Controllers
func InitManagedUserController() identitycontrollers.UserSignInController {
	wire.Build(Container)
	return identitycontrollers.UserSignInController{}
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

func InitResetPasswordController() identitycontrollers.ResetPasswordController {
	wire.Build(Container)
	return identitycontrollers.ResetPasswordController{}
}

// Event Handlers

func InitVerifcationEmailSender() identityemail.VerifcationEmailSender {
	wire.Build(Container)
	return identityemail.VerifcationEmailSender{}
}
