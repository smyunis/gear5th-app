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

var editAdSlotTemplate *template.Template

func init() {
	editAdSlotTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/ads/edit-adslot.html"))
}

type editAdSlotPresenter struct {
	AdSlotName           string `form:"adslot-name"`
	AdSlotType           string
	AdSlotTypeDimentions string
	ErrorMessage         string
}

type EditAdSlotController struct {
	authMiddleware   middlewares.JwtAuthenticationMiddleware
	adSlotInteractor adslotinteractors.AdSlotInteractor
	logger           application.Logger
}

func NewEditAdSlotController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	adSlotsInteractor adslotinteractors.AdSlotInteractor,
	logger application.Logger) EditAdSlotController {
	return EditAdSlotController{
		authMiddleware,
		adSlotsInteractor,
		logger,
	}
}

func (c *EditAdSlotController) AddRoutes(router *fiber.Router) {
	(*router).Use("/ads/edit-adslot", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/ads/edit-adslot/:adSlotID", c.editAdSlotOnGet)
	(*router).Add(fiber.MethodPost, "/ads/edit-adslot/:adSlotID", c.editAdSlotOnPost)
}

func (c *EditAdSlotController) editAdSlotOnGet(ctx *fiber.Ctx) error {
	adSlotID := ctx.Params("adSlotID", "")
	if adSlotID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	slot, err := c.adSlotInteractor.AdSlot(shared.ID(adSlotID))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	p := editAdSlotPresenter{
		AdSlotName:           slot.Name(),
		AdSlotType:           slot.AdSlotType().String(),
		AdSlotTypeDimentions: slot.AdSlotType().Dimentions().String(),
	}
	return controllers.Render(ctx, editAdSlotTemplate, p)
}

func (c *EditAdSlotController) editAdSlotOnPost(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	adSlotID := ctx.Params("adSlotID", "")
	if adSlotID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	p := &editAdSlotPresenter{
	}
	err = ctx.BodyParser(p)
	if err != nil {
		p.ErrorMessage = "One or more invalid inputs. Check and try again."
		return controllers.Render(ctx, editAdSlotTemplate, p)
	}

	err = c.adSlotInteractor.ChangeAdSlotName(publisherID, shared.ID(adSlotID), p.AdSlotName)
	if err != nil {
		c.logger.Error("adslots/changename", err)
		p.ErrorMessage = "We're unable to change the name of your ad slot at the moment. Try again later."
		return controllers.Render(ctx, editAdSlotTemplate, p)
	}

	return ctx.Redirect("/publish/ads")
}
