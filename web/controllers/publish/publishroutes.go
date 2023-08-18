package publish

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/web/ioc"
)

func Routes(app *fiber.App) {

	publishRouter := app.Group("/publish")

	managedUserController := ioc.InitManagedUserController()
	// publishRouter.Add(fiber.MethodGet, "/identity/signin", managedUserController.Get)
	// publishRouter.Add(fiber.MethodPost, "/identity/signin", managedUserController.Post)
	managedUserController.AddRoutes(&publishRouter)

	publisherSignUpController := ioc.InitPublisherSignUpController()
	publishRouter.Add(publisherSignUpController.Method, publisherSignUpController.Path, publisherSignUpController.ManagedUserSignUp)

	requestPasswordResetController := ioc.InitRequestPasswordResetController()
	publishRouter.Add(requestPasswordResetController.Method, requestPasswordResetController.Path, requestPasswordResetController.RequestPasswordReset)

	verifyEmailController := ioc.InitVerifyEmailController()
	publishRouter.Add(verifyEmailController.Method, verifyEmailController.Path, verifyEmailController.VerifyEmail)

	resetPasswordController := ioc.InitResetPasswordController()
	publishRouter.Add(resetPasswordController.Method, resetPasswordController.Path, resetPasswordController.ResetPassword)

}
