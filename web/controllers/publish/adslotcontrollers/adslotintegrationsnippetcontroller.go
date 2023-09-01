package adslotcontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adslotinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var adSlotIntegrationSnippetTemplate *template.Template

func init() {
	adSlotIntegrationSnippetTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/ads/integration-snippet.html"))
}

type adSlotIntegrationSnippetPresenter struct {
	Nav         string
	HTMLSnippet string
}

type AdSlotIntegrationSnippetController struct {
	authMiddleware   middlewares.JwtAuthenticationMiddleware
	adSlotInteractor adslotinteractors.AdSlotInteractor
	logger           application.Logger
}

func NewAdSlotIntegrationSnippetController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	adSlotsInteractor adslotinteractors.AdSlotInteractor,
	logger application.Logger) AdSlotIntegrationSnippetController {
	return AdSlotIntegrationSnippetController{
		authMiddleware,
		adSlotsInteractor,
		logger,
	}
}

func (c *AdSlotIntegrationSnippetController) AddRoutes(router *fiber.Router) {
	(*router).Use("/ads/integration-snippet", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/ads/integration-snippet/:adSlotID", c.integrationSnippetOnGet)
}

func (c *AdSlotIntegrationSnippetController) integrationSnippetOnGet(ctx *fiber.Ctx) error {
	adSlotID := ctx.Params("adSlotID", "")
	if adSlotID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	snippet, err := c.adSlotInteractor.GetIntegrationHTMLSnippet(shared.ID(adSlotID))
	if err != nil {
		c.logger.Error("adslots/integrationsnippet", err)
		return ctx.Redirect("/pages/error.html")
	}

	p := adSlotIntegrationSnippetPresenter{
		Nav:         "ads",
		HTMLSnippet: snippet,
	}

	return controllers.Render(ctx, adSlotIntegrationSnippetTemplate, p)

}
