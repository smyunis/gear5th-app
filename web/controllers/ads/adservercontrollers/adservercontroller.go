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

type adServerPresenter struct {
	SiteID      string
	PublisherID string
	AdSlotID    string
	AdView      adsinteractors.AdView
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

	siteID := ctx.Query("site-id", "")
	if siteID == "" {
		return nil
	}

	slotID := ctx.Query("adslot-id", "")
	if slotID == "" {
		return nil
	}

	publisherID := ctx.Query("publisher-id", "")
	if publisherID == "" {
		return nil
	}


	ad, err := c.adsPool.Next(adslot.AdSlotTypeFromString(slot))
	if err != nil {
		c.logger.Error("ads/adserver", err)
		return nil
	}

	p := adServerPresenter{
		SiteID:      siteID,
		AdSlotID:    slotID,
		PublisherID: publisherID,
		AdView:      *ad,
	}

	return controllers.Render(ctx, adServerTemplate, p)
}
