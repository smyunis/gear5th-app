package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofor-little/env"
	"gitlab.com/gear5th/gear5th-api/cmd/api/ioc"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	err := env.Load("config/.env.dev")
	if err != nil {
		panic("could not load config file ./config/.env.dev")
	}

	app := fiber.New()
	app.Use(recover.New())
	addIdentityRoutes(app)

	app.Listen(":5071")
}

func addIdentityRoutes(app *fiber.App) {
	identityRouter := app.Group("/identity")
	identityRouter.Post("/managed/signin", ioc.InitManagedUserSignInController().SignIn)
}
