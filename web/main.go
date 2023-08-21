package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"gitlab.com/gear5th/gear5th-api/web/controllers/publish"
)

func main() {

	err := godotenv.Load("config/.env.dev", "config/.env.prod")
	if err != nil {
		panic("could not load config file ./config/.env.*")
	}

	registerEventHandlers()

	engine := html.New("./web/views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(recover.New())
	app.Static("/", "web/public")
	publish.AddRoutes(app)

	app.Listen(":5071")
}

func registerEventHandlers() {
}
