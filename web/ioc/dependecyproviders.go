//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/adslotcontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/homecontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/identitycontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/publishercontrollers"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/sitecontrollers"
	"gitlab.com/gear5th/gear5th-app/web/events"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

// Middlewares

func InitJwtAuthenticationMiddleware() middlewares.JwtAuthenticationMiddleware {
	wire.Build(Container)
	return middlewares.JwtAuthenticationMiddleware{}
}

// Controllers
func InitManagedUserController() identitycontrollers.UserSignInController {
	wire.Build(Container)
	return identitycontrollers.UserSignInController{}
}

func InitOAuthSignInController() identitycontrollers.OAuthSignInController {
	wire.Build(Container)
	return identitycontrollers.OAuthSignInController{}
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

func InitHomeController() homecontrollers.HomeController {
	wire.Build(Container)
	return homecontrollers.HomeController{}
}

func InitSiteController() sitecontrollers.SiteController {
	wire.Build(Container)
	return sitecontrollers.SiteController{}
}

func InitCreateSiteController() sitecontrollers.CreateSiteController {
	wire.Build(Container)
	return sitecontrollers.CreateSiteController{}
}

func InitVerifySiteController() sitecontrollers.VerifySiteController {
	wire.Build(Container)
	return sitecontrollers.VerifySiteController{}
}

func InitAdSlotController() adslotcontrollers.AdSlotController {
	wire.Build(Container)
	return adslotcontrollers.AdSlotController{}
}

func InitCreateAdSlotController() adslotcontrollers.CreateAdSlotController {
	wire.Build(Container)
	return adslotcontrollers.CreateAdSlotController{}
}

func InitEditAdSlotController() adslotcontrollers.EditAdSlotController {
	wire.Build(Container)
	return adslotcontrollers.EditAdSlotController{}
}

func InitAdSlotIntegrationSnippetController() adslotcontrollers.AdSlotIntegrationSnippetController {
	wire.Build(Container)
	return adslotcontrollers.AdSlotIntegrationSnippetController{}
}

// Event Handlers

func InitEventsRegistrar() events.EventHandlerRegistrar {
	wire.Build(Container)
	return events.EventHandlerRegistrar{}
}
