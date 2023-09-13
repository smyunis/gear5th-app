package campaigncontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/advertiserinteractors"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var campaignTemplate *template.Template

func init() {
	campaignTemplate = template.Must(
		controllers.AdvertiserMainLayoutTemplate().ParseFiles(
			"web/views/advertiser/campaign/campaign.html"))

}

type advertiserCampigns struct {
	ID        string
	Name      string
	Start     string
	End       string
	IsRunning bool
}

type campaignPresenter struct {
	Campaigns       []advertiserCampigns
	Token           string
	AdvertiserEmail string
}

type CampaignController struct {
	advertiserRefferal middlewares.AdvertiserRefferalMiddleware
	campaignInteractor advertiserinteractors.CampaignInteractor
	store              application.FileStore
	logger             application.Logger
}

func NewCampaignController(
	advertiserRefferal middlewares.AdvertiserRefferalMiddleware,
	campaignInteractor advertiserinteractors.CampaignInteractor,
	store application.FileStore,
	logger application.Logger) CampaignController {
	return CampaignController{
		advertiserRefferal,
		campaignInteractor,
		store,
		logger,
	}
}

func (c *CampaignController) AddRoutes(router *fiber.Router) {
	(*router).Use("/campaign", c.advertiserRefferal.Verification)
	(*router).Add(fiber.MethodGet, "/campaign", c.campaignOnGet)
}

// 01HA4GNAXY6DKN506Z4HXM8DTY
// 01HA4GNAXY6DKN506Z4HXM8DTY b5b3323c6b96ad10ce4fe489b950b2b050eebf9d4fefbb3e3ac49404e112a83e
// http://localhost:5071/advertiser/campaign?token=01HA4GNAXY6DKN506Z4HXM8DTY%20b5b3323c6b96ad10ce4fe489b950b2b050eebf9d4fefbb3e3ac49404e112a83e
func (c *CampaignController) campaignOnGet(ctx *fiber.Ctx) error {
	advertiserID, err := c.advertiserRefferal.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	token, err := c.advertiserRefferal.AdvertiserToken(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	ad, err := c.campaignInteractor.Advertiser(advertiserID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	campaigns, err := c.campaignInteractor.CampaignsForAdvertiser(advertiserID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	advertiserCamps := make([]advertiserCampigns, 0)
	for _, camp := range campaigns {
		advertiserCamps = append(advertiserCamps, advertiserCampigns{
			ID:        camp.ID.String(),
			Name:      camp.Name,
			Start:     camp.Start.Format("02-Jan-2006"),
			End:       camp.End.Format("02-Jan-2006"),
			IsRunning: camp.IsRunning(),
		})
	}

	p := &campaignPresenter{
		advertiserCamps,
		token,
		ad.Email().String(),
	}

	return controllers.Render(ctx, campaignTemplate, p)
}
