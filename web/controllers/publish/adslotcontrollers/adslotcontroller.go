package adslotcontrollers

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var adsTemplate *template.Template

func init() {
	adsTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/ads/ads.html"))
}

type adSlotsForSitePresenter struct {
	AdSlotID       string
	Name           string
	AdSlotType     string
	SlotDimentions string
}

type adSlotPresenter struct {
	ErrorMessage string
	AdSlots      map[string][]adSlotsForSitePresenter
}

type AdSlotController struct {
	authMiddleware   middlewares.JwtAuthenticationMiddleware
	adSlotInteractor publisherinteractors.AdSlotInteractor
	logger           application.Logger
}

func NewAdSlotController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	adSlotsInteractor publisherinteractors.AdSlotInteractor,
	logger application.Logger) AdSlotController {
	return AdSlotController{
		authMiddleware,
		adSlotsInteractor,
		logger,
	}
}

func (c *AdSlotController) AddRoutes(router *fiber.Router) {
	(*router).Use("/ads", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/ads", c.adsPageOnGet)
	(*router).Add(fiber.MethodGet, "/ads/remove-adslot/:adSlotID", c.removeAdSlot)
}

func (c *AdSlotController) adsPageOnGet(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	slots, err := c.adSlotInteractor.ActiveAdSlotsForPublisher(publisherID)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			p := adSlotPresenter{
				ErrorMessage: "There are no ad slots registred for any of your sites",
			}
			return controllers.Render(ctx, adsTemplate, p)
		}
		p := adSlotPresenter{
			ErrorMessage: "We're unable to get your ad slots for site at the moment. Try again later.",
		}
		c.logger.Error("adslot/activesitesforpublisher", err)
		return controllers.Render(ctx, adsTemplate, p)
	}

	adSlots := make(map[string][]adSlotsForSitePresenter)
	for siteDomain, siteSlots := range slots {
		adSlotsForSite := make([]adSlotsForSitePresenter, 0)
		for _, slot := range siteSlots {
			adSlotsForSite = append(adSlotsForSite, adSlotsForSitePresenter{
				AdSlotID:       slot.ID.String(),
				Name:           slot.Name,
				AdSlotType:     slot.SlotType.String(),
				SlotDimentions: fmt.Sprintf("%dx%d", slot.SlotType.Dimentions().Width, slot.SlotType.Dimentions().Height),
			})
		}
		adSlots[siteDomain] = adSlotsForSite
	}

	p := adSlotPresenter{
		AdSlots: adSlots,
	}
	return controllers.Render(ctx, adsTemplate, p)
}

func (c *AdSlotController) removeAdSlot(ctx *fiber.Ctx) error {
	siteID := ctx.Params("adSlotID", "")
	if siteID == "" {
		return ctx.Redirect("/pages/error.html")
	}
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	err = c.adSlotInteractor.DeactivateAdSlot(publisherID, shared.ID(siteID))
	if err != nil {
		c.logger.Error("adslot/removeadslot", err)
		return ctx.Redirect("/pages/error.html")
	}

	return ctx.Redirect("/publish/ads")
}
