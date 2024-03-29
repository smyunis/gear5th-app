package adpiececontrollers

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/advertiserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/adslot"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var addAdPiecesTemplate *template.Template

func init() {
	addAdPiecesTemplate = template.Must(
		controllers.AdvertiserMainLayoutTemplate().ParseFiles(
			"web/views/advertiser/adpiece/add-adpiece.html"))

}

type addAdPiecesPresenter struct {
	Ref          string `form:"ref"`
	SlotType     string `form:"adslot-type"`
	ErrorMessage string
	Campaign     campaign.Campaign
	Token        string
}

type AddAdPieceController struct {
	advertiserRefferal middlewares.AdvertiserRefferalMiddleware
	adPieceInteractor  advertiserinteractors.AdPieceInteractor
	campaignInteractor advertiserinteractors.CampaignInteractor
	store              application.FileStore
	logger             application.Logger
}

func NewAddAdPieceController(
	advertiserRefferal middlewares.AdvertiserRefferalMiddleware,
	adPieceInteractor advertiserinteractors.AdPieceInteractor,
	campaignInteractor advertiserinteractors.CampaignInteractor,
	store application.FileStore,
	logger application.Logger) AddAdPieceController {
	return AddAdPieceController{
		advertiserRefferal,
		adPieceInteractor,
		campaignInteractor,
		store,
		logger,
	}
}

func (c *AddAdPieceController) AddRoutes(router *fiber.Router) {
	(*router).Use("/adpiece", c.advertiserRefferal.Verification)
	(*router).Add(fiber.MethodGet, "/adpiece/add-adpiece", c.addAdPiecesOnGet)
	(*router).Add(fiber.MethodPost, "/adpiece/add-adpiece", c.addAdPiecesOnPost)
}

// http://localhost:5071/advertiser/adpiece/add-adpiece?campaignId=01HA4GNAZA5D2V60SEDDS28TRM&token=01HA4GNAXY6DKN506Z4HXM8DTY%20b5b3323c6b96ad10ce4fe489b950b2b050eebf9d4fefbb3e3ac49404e112a83e
func (c *AddAdPieceController) addAdPiecesOnGet(ctx *fiber.Ctx) error {
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
		c.logger.Error("adpiece/addadpiece/getcampaign", err)
		return ctx.Redirect("/pages/error.html")
	}

	p := &addAdPiecesPresenter{
		Token:    token,
		Campaign: camp,
	}
	return controllers.Render(ctx, addAdPiecesTemplate, p)
}

func (c *AddAdPieceController) addAdPiecesOnPost(ctx *fiber.Ctx) error {
	actorID, err := c.advertiserRefferal.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	campaignID := ctx.Query("campaignId", "")
	if campaignID == "" {
		return ctx.Redirect("/pages/error.html")
	}
	token, err := c.advertiserRefferal.AdvertiserToken(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	p := &addAdPiecesPresenter{}
	err = ctx.BodyParser(p)
	if err != nil {
		p.ErrorMessage = "One or more invalid inputs. Check and try again"
		return controllers.Render(ctx, addAdPiecesTemplate, p)
	}

	resource, err := ctx.FormFile("resource")
	if err != nil {
		p.ErrorMessage = "Unable to upload resource image. Check and try again."
		return controllers.Render(ctx, addAdPiecesTemplate, p)
	}
	resourceMIME := resource.Header.Get("Content-Type")
	resourceReader, err := resource.Open()
	if err != nil {
		p.ErrorMessage = "Unable to upload resource image. Check and try again."
		return controllers.Render(ctx, addAdPiecesTemplate, p)
	}

	refLink, err := url.Parse(p.Ref)
	if err != nil {
		refLink, _ = url.Parse("#")
	}
	err = c.campaignInteractor.AddAdPiece(
		actorID, shared.ID(campaignID), adSlotType(p.SlotType), refLink, resourceMIME, resourceReader)

	if err != nil {
		c.logger.Error("adpiece/addadpiece", err)
		p.ErrorMessage = "We're unable to add a new adpiece at at the moment. Try again later."
		return controllers.Render(ctx, addAdPiecesTemplate, p)
	}

	return ctx.Redirect(fmt.Sprintf("/advertiser/adpiece?campaignId=%s&token=%s", campaignID, token))
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
