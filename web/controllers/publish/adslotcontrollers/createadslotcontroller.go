package adslotcontrollers

import (
	"errors"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adslotinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/siteinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var createAdSlotTemplate *template.Template

func init() {
	createAdSlotTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/ads/create-adslot.html"))
}

type createAdSlotSitesPresenter struct {
	SiteDomain string
	SiteID     string
}

type createAdSlotPresenter struct {
	AdSlotName   string `form:"adslot-name"`
	SlotType     string `form:"adslot-type"`
	SiteID       string `form:"adslot-site-id"`
	ErrorMessage string
	Sites        []createAdSlotSitesPresenter
}

type CreateAdSlotController struct {
	authMiddleware   middlewares.JwtAuthenticationMiddleware
	adSlotInteractor adslotinteractors.AdSlotInteractor
	sitesInteractor  siteinteractors.SiteInteractor
	logger           application.Logger
}

func NewCreateAdSlotController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	adSlotsInteractor adslotinteractors.AdSlotInteractor,
	sitesInteractor siteinteractors.SiteInteractor,
	logger application.Logger) CreateAdSlotController {
	return CreateAdSlotController{
		authMiddleware,
		adSlotsInteractor,
		sitesInteractor,
		logger,
	}
}

func (c *CreateAdSlotController) AddRoutes(router *fiber.Router) {
	(*router).Use("/ads/create-adslot", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/ads/create-adslot", c.createAdSlotOnGet)
	(*router).Add(fiber.MethodPost, "/ads/create-adslot", c.createAdSlotOnPost)
}

func (c *CreateAdSlotController) createAdSlotOnGet(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	activeSites, err := c.sitesInteractor.ActiveSitesForPublisher(publisherID)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			p := createAdSlotPresenter{
				ErrorMessage: "You have no sites registered. You must register a site first so you can add ad slots to it.",
			}
			return controllers.Render(ctx, createAdSlotTemplate, p)
		}

		c.logger.Error("adslots/createadslot", err)
		return ctx.Redirect("/pages/error.html")
	}

	sites := make([]createAdSlotSitesPresenter, 0)
	for _, s := range activeSites {
		sites = append(sites, createAdSlotSitesPresenter{
			SiteID:     s.ID().String(),
			SiteDomain: s.SiteDomain(),
		})
	}

	p := createAdSlotPresenter{
		Sites: sites,
	}
	return controllers.Render(ctx, createAdSlotTemplate, p)
}

func (c *CreateAdSlotController) createAdSlotOnPost(ctx *fiber.Ctx) error {
	p := &createAdSlotPresenter{}
	err := ctx.BodyParser(p)
	if err != nil {
		c.logger.Error("adslots/createadslot", err)
		return ctx.Redirect("/pages/error.html")
	}
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	
	err = c.adSlotInteractor.CreateAdSlotForSite(publisherID, shared.ID(p.SiteID), p.AdSlotName, adSlotType(p.SlotType))
	if err != nil {
		p.ErrorMessage = "We're unable to create your ad slot at the moment. Try agian later."
		c.logger.Error("adslots/createadslot", err)
		return controllers.Render(ctx, createAdSlotTemplate, p)
	}
	return ctx.Redirect("/publish/ads")
}

func adSlotType(slotType string) adslot.AdSlotType {
	switch slotType {
	case "horizontal":
		return adslot.Horizontal
	case "vertical":
		return adslot.Vertical
	case "box":
		return adslot.Box
	default:
		return 0
	}
}
