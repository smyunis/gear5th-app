package adservercontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

var adServerTemplate *template.Template

func init() {
	adServerTemplate = template.Must(
		template.ParseFiles(
			"web/views/ads/adserver.html"))

}

type AdServerController struct {
	adsPool adsinteractors.AdsPool
	logger  application.Logger
}

func NewAdServerController(
	adsPool adsinteractors.AdsPool,
	logger application.Logger) AdServerController {
	return AdServerController{
		adsPool,
		logger,
	}
}

func (c *AdServerController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/adserver", c.adServerOnGet)
}

func (c *AdServerController) adServerOnGet(ctx *fiber.Ctx) error {
	slot := ctx.Query("slot", "")
	if slot == "" {
		return nil
	}

	ad, err := c.adsPool.Next(adslot.AdSlotTypeFromString(slot))
	if err != nil {
		c.logger.Error("ads/adserver", err)
		return nil
	}

	return controllers.Render(ctx, adServerTemplate, ad)
}
