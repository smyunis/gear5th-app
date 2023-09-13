package advertiser

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/web/ioc"
)

func Routes(app *fiber.App) {

	route := app.Group("/advertiser")

	adPieceController := ioc.InitAdPieceController()
	adPieceController.AddRoutes(&route)

	campaignController := ioc.InitCampaignController()
	campaignController.AddRoutes(&route)

}
