package admin

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/web/ioc"
)

func Routes(app *fiber.App) {

	route := app.Group("/admin")

	identityController := ioc.InitAdminIdentityController()
	identityController.AddRoutes(&route)

	dashController := ioc.InitAdminDashboardController()
	dashController.AddRoutes(&route)

	pubPay := ioc.InitAdminPublisherPaymentsController()
	pubPay.AddRoutes(&route)

	adv := ioc.InitAdminAdvertisersController()
	adv.AddRoutes(&route)

}
