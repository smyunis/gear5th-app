package publish

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish/homecontrollers"
	"gitlab.com/gear5th/gear5th-app/web/ioc"
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
	requestPasswordResetController.AddRoutes(&publishRouter)

	verifyEmailController := ioc.InitVerifyEmailController()
	verifyEmailController.AddRoutes(&publishRouter)

	resetPasswordController := ioc.InitResetPasswordController()
	resetPasswordController.AddRoutes(&publishRouter)
}
