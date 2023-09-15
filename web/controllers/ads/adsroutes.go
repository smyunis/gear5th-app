package ads

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/web/ioc"
)

func Routes(app *fiber.App) {

	route := app.Group("/ads")

	adServerController := ioc.InitAdServerController()
	adServerController.AddRoutes(&route)
}
