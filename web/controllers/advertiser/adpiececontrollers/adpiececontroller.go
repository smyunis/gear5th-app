package adpiececontrollers

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/advertiserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/adpiece"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var adPiecesTemplate *template.Template

func init() {
	adPiecesTemplate = template.Must(
		controllers.AdvertiserMainLayoutTemplate().ParseFiles(
			"web/views/advertiser/adpiece/adpiece.html"))

}

type adPiecesPresenter struct {
	AdPieces []adpiece.AdPiece
	Campaign campaign.Campaign
	Token    string
}

type AdPieceController struct {
	advertiserRefferal middlewares.AdvertiserRefferalMiddleware
	adPieceInteractor  advertiserinteractors.AdPieceInteractor
	campaignInteractor advertiserinteractors.CampaignInteractor
	store              application.FileStore
	logger             application.Logger
}

func NewAdPieceController(
	advertiserRefferal middlewares.AdvertiserRefferalMiddleware,
	adPieceInteractor advertiserinteractors.AdPieceInteractor,
	campaignInteractor advertiserinteractors.CampaignInteractor,
	store application.FileStore,
	logger application.Logger) AdPieceController {
	return AdPieceController{
		advertiserRefferal,
		adPieceInteractor,
		campaignInteractor,
		store,
		logger,
	}
}

func (c *AdPieceController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/adpiece/:adPieceId/resource", c.resourceOnGet)
	
	(*router).Use("/adpiece", c.advertiserRefferal.Verification)
	(*router).Add(fiber.MethodGet, "/adpiece", c.adPiecesOnGet)
	(*router).Add(fiber.MethodGet, "/adpiece/:adPieceId/remove", c.removeAdPieceOnGet)
}

func (c *AdPieceController) adPiecesOnGet(ctx *fiber.Ctx) error {
	campaignID := ctx.Query("campaignId", "")
	if campaignID == "" {
		return ctx.Redirect("/pages/error.html")
	}
	token, err := c.advertiserRefferal.AdvertiserToken(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	camp, err := c.campaignInteractor.Campaign(shared.ID(campaignID))
	if err != nil {
		c.logger.Error("adpiece/campaign/get", err)
		return ctx.Redirect("/pages/error.html")
	}

	adPieces, err := c.adPieceInteractor.ActiveAdPiecesForCampaign(shared.ID(campaignID))
	if err != nil && !errors.Is(err, application.ErrEntityNotFound) {
		c.logger.Error("adpiece/campaign/activeadpieces", err)
		return ctx.Redirect("/pages/error.html")
	}

	p := &adPiecesPresenter{
		adPieces,
		camp,
		token,
	}

	return controllers.Render(ctx, adPiecesTemplate, p)
}

// http://localhost:5071/advertiser/adpiece/01HA4VB89SS6BW4BCQYXTZ1R21/resource
func (c *AdPieceController) resourceOnGet(ctx *fiber.Ctx) error {
	//TODO Aggressively cache this endpoint, its doom of it gets hit for every request
	adPieceID := ctx.Params("adPieceId", "")
	if adPieceID == "" {
		return nil
	}
	a, err := c.adPieceInteractor.AdPiece(shared.ID(adPieceID))
	if err != nil {
		c.logger.Error("adpiece/get", err)
		return nil
	}

	resource, err := c.store.Get(a.Resource)
	if err != nil {
		c.logger.Error("adpiece/get/resource", err)
		return nil
	}

	ctx.Set("Content-Type", a.ResourceMIMEType)
	res := ctx.Response()
	res.StreamBody = true
	res.SetBodyStream(resource, -1)

	return nil
}



func (c *AdPieceController) removeAdPieceOnGet(ctx *fiber.Ctx) error {
	campaignID := ctx.Query("campaignId", "")
	if campaignID == "" {
		return ctx.Redirect("/pages/error.html")
	}
	token, err := c.advertiserRefferal.AdvertiserToken(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	actorID, err := c.advertiserRefferal.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	adPieceID := ctx.Params("adPieceId", "")
	if adPieceID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	err = c.adPieceInteractor.DeactivateAdPiece(actorID, shared.ID(adPieceID))
	if err != nil {
		c.logger.Error("adpiece/deactivate", err)
		return ctx.Redirect("/pages/error.html")
	}

	return ctx.Redirect(fmt.Sprintf("/advertiser/adpiece?campaignId=%s&token=%s", campaignID, token))
}
