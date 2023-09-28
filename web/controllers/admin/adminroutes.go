package admin

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {

	_ = app.Group("/admin")

}
