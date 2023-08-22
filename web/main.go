package main

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gitlab.com/gear5th/gear5th-app/web/controllers/publish"
)

func main() {

	err := godotenv.Load("config/.env.dev", "config/.env.prod")
	if err != nil {
		panic("could not load config file ./config/.env.*")
	}

	registerEventHandlers()

	app := fiber.New()

	app.Use(recover.New())
	app.Static("/", "web/public")
	publish.Routes(app)

	app.Listen(":5071")
}

func registerEventHandlers() {
}
