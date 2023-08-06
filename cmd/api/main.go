package main

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/cmd/api/dependencies"
)

func main() {
	app := fiber.New()

	addIdentityRoutes(app)

	app.Listen(":5071")
}

func addIdentityRoutes(app *fiber.App) {
	identityRouter := app.Group("/identity")
	identityRouter.Post("/managed/signin", dependencies.InitManagedUserSignInController().SignIn)
}
