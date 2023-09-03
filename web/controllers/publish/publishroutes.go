package publish

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/web/ioc"
)

func Routes(app *fiber.App) {

	publishRouter := app.Group("/publish")

	managedUserController := ioc.InitManagedUserController()
	managedUserController.AddRoutes(&publishRouter)

	publisherSignUpController := ioc.InitPublisherSignUpController()
	publisherSignUpController.AddRoutes(&publishRouter)

	homeController := ioc.InitHomeController()
	homeController.AddRoutes(&publishRouter)

	requestPasswordResetController := ioc.InitRequestPasswordResetController()
	requestPasswordResetController.AddRoutes(&publishRouter)

	verifyEmailController := ioc.InitVerifyEmailController()
	verifyEmailController.AddRoutes(&publishRouter)

	resetPasswordController := ioc.InitResetPasswordController()
	resetPasswordController.AddRoutes(&publishRouter)

	siteController := ioc.InitSiteController()
	siteController.AddRoutes(&publishRouter)

	createSiteController := ioc.InitCreateSiteController()
	createSiteController.AddRoutes(&publishRouter)

	verifySiteController := ioc.InitVerifySiteController()
	verifySiteController.AddRoutes(&publishRouter)

	adSlotController := ioc.InitAdSlotController()
	adSlotController.AddRoutes(&publishRouter)

	createAdSlotController := ioc.InitCreateAdSlotController()
	createAdSlotController.AddRoutes(&publishRouter)

	editAdSlotController := ioc.InitEditAdSlotController()
	editAdSlotController.AddRoutes(&publishRouter)

	adSlotIntegrationSnippetController := ioc.InitAdSlotIntegrationSnippetController()
	adSlotIntegrationSnippetController.AddRoutes(&publishRouter)

	oauthUserSigninController := ioc.InitOAuthSignInController()
	oauthUserSigninController.AddRoutes(&publishRouter)

	accountcontroller := ioc.InitAccountController()
	accountcontroller.AddRoutes(&publishRouter)

}
