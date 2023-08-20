package publish

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/web/controllers/publish/homecontrollers"
	"gitlab.com/gear5th/gear5th-api/web/ioc"
)

func Routes(app *fiber.App) {

	publishRouter := app.Group("/publish")

	managedUserController := ioc.InitManagedUserController()
	managedUserController.AddRoutes(&publishRouter)

	publisherSignUpController := ioc.InitPublisherSignUpController()
	publisherSignUpController.AddRoutes(&publishRouter)

	//TODO get from DI
	homeController := homecontrollers.NewHomeController()
	homeController.AddRoutes(&publishRouter)

	requestPasswordResetController := ioc.InitRequestPasswordResetController()
	publishRouter.Add(requestPasswordResetController.Method, requestPasswordResetController.Path, requestPasswordResetController.RequestPasswordReset)

	verifyEmailController := ioc.InitVerifyEmailController()
	publishRouter.Add(verifyEmailController.Method, verifyEmailController.Path, verifyEmailController.VerifyEmail)

	resetPasswordController := ioc.InitResetPasswordController()
	publishRouter.Add(resetPasswordController.Method, resetPasswordController.Path, resetPasswordController.ResetPassword)

}
