package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gitlab.com/gear5th/gear5th-api/cmd/api/ioc"

	// Added to register domain event handlers in their init functions
	"gitlab.com/gear5th/gear5th-api/internal/application"
	_ "gitlab.com/gear5th/gear5th-api/internal/infrastructure/mail/identityemail"
)

func main() {

	err := godotenv.Load("config/.env.dev", "config/.env.prod")
	if err != nil {
		panic("could not load config file ./config/.env.*")
	}

	registerEventHandlers()

	app := fiber.New()
	app.Use(recover.New())
	addRoutes(app)

	app.Listen(":5071")
}

func addRoutes(app *fiber.App) {
	identityRouter := app.Group("/")

	managedUserController := ioc.InitManagedUserController()
	identityRouter.Add(managedUserController.Method, managedUserController.Path, managedUserController.SignIn)

	publisherSignUpController := ioc.InitPublisherSignUpController()
	identityRouter.Add(publisherSignUpController.Method, publisherSignUpController.Path, publisherSignUpController.ManagedUserSignUp)

	requestPasswordResetController := ioc.InitRequestPasswordResetController()
	identityRouter.Add(requestPasswordResetController.Method, requestPasswordResetController.Path, requestPasswordResetController.RequestPasswordReset)

}

func registerEventHandlers() {
	emailVerificationSender := ioc.InitVerifcationEmailSender()
	application.ApplicationEventDispatcher.AddHandler("user.signedup", emailVerificationSender.SendMail)
}
