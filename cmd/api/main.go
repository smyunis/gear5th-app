package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofor-little/env"
	"gitlab.com/gear5th/gear5th-api/cmd/api/ioc"
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
	identityRouter := app.Group("/")

	controller := ioc.InitManagedUserController()
	identityRouter.Add(controller.Method, controller.Path, controller.SignIn)
}
