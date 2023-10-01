package adminadvertisers

import (
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/advertiserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/advertiser/campaign"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/deposit"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var adminAdvertisersTemplate *template.Template
var adminAdvertiserTemplate *template.Template
var adminNewAdvertiserTemplate *template.Template

func init() {
	adminAdvertisersTemplate = template.Must(
		controllers.AdminMainLayoutTemplate().ParseFiles(
			"web/views/admin/advertisers/advertisers.html"))
	adminAdvertiserTemplate = template.Must(
		controllers.AdminMainLayoutTemplate().ParseFiles(
			"web/views/admin/advertisers/advertiser.html"))
	adminNewAdvertiserTemplate = template.Must(
		controllers.AdminMainLayoutTemplate().ParseFiles(
			"web/views/admin/advertisers/new-advertiser.html"))
}

type adminAdvertisersPresenter struct {
}

type AdminAdvertisersController struct {
	authMiddleware       middlewares.AdminAuthenticationMiddleware
	advertiserInteractor advertiserinteractors.AdvertiserInteractor
	campaignInteractor   advertiserinteractors.CampaignInteractor
	depositInteractor    paymentinteractors.DepositInteractor
	logger               application.Logger
}

func NewAdminAdvertisersController(
	authMiddleware middlewares.AdminAuthenticationMiddleware,
	advertiserInteractor advertiserinteractors.AdvertiserInteractor,
	campaignInteractor advertiserinteractors.CampaignInteractor,
	depositInteractor paymentinteractors.DepositInteractor,
	logger application.Logger) AdminAdvertisersController {
	return AdminAdvertisersController{
		authMiddleware,
		advertiserInteractor,
		campaignInteractor,
		depositInteractor,
		logger,
	}
}

func (c *AdminAdvertisersController) AddRoutes(router *fiber.Router) {
	(*router).Use("/advertisers", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/advertisers", c.advertisersOnGet)
	(*router).Add(fiber.MethodGet, "/advertisers/new-advertiser", c.newAdvertiserOnGet)
	(*router).Add(fiber.MethodPost, "/advertisers/new-advertiser", c.newAdvertiserOnPost)
	(*router).Add(fiber.MethodPost, "/advertisers/:advertiserId/campaign/:campaignId/quit", c.quitCampaignOnPost)
	(*router).Add(fiber.MethodPost, "/advertisers/:advertiserId/campaign", c.newCampaignOnPost)
	(*router).Add(fiber.MethodPost, "/advertisers/:advertiserId/deposit", c.acceptDepositOnPost)
	(*router).Add(fiber.MethodGet, "/advertisers/:advertiserId", c.advertiserOnGet)
}

func (c *AdminAdvertisersController) advertisersOnGet(ctx *fiber.Ctx) error {

	advertisers, err := c.advertiserInteractor.Adveritsers()
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	return controllers.Render(ctx, adminAdvertisersTemplate, advertisers)
}

type advertieserDetailsPresenter struct {
	advertiserinteractors.AdveritserDetails
	Deposits  []deposit.Deposit
	Campaigns []campaign.Campaign
	Token     string
}

func (c *AdminAdvertisersController) advertiserOnGet(ctx *fiber.Ctx) error {
	advertiserID := ctx.Params("advertiserId", "")
	if advertiserID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	p := advertieserDetailsPresenter{}

	ad, err := c.advertiserInteractor.Advertiser(shared.ID(advertiserID))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	p.AdveritserDetails = ad

	deps, err := c.depositInteractor.DepoistsForAdvertiser(shared.ID(advertiserID))
	if err != nil {
		p.Deposits = make([]deposit.Deposit, 0)
		return controllers.Render(ctx, adminAdvertiserTemplate, p)
	}
	p.Deposits = deps

	campaigns, err := c.campaignInteractor.CampaignsForAdvertiser(shared.ID(advertiserID))
	if err != nil {
		p.Campaigns = make([]campaign.Campaign, 0)
		return controllers.Render(ctx, adminAdvertiserTemplate, p)
	}

	p.Campaigns = campaigns

	token, err := c.advertiserInteractor.GenerateAdvertiserToken(shared.ID(advertiserID))
	if err != nil {
		return controllers.Render(ctx, adminAdvertiserTemplate, p)
	}

	p.Token = token

	return controllers.Render(ctx, adminAdvertiserTemplate, p)
}

func (c *AdminAdvertisersController) newAdvertiserOnGet(ctx *fiber.Ctx) error {
	return controllers.Render(ctx, adminNewAdvertiserTemplate, nil)
}

type newAdvertiserPresenter struct {
	Name         string `form:"name"`
	Email        string `form:"email"`
	PhoneNumber  string `form:"phone-number"`
	Note         string `form:"note"`
	ErrorMessage string
}

func (c *AdminAdvertisersController) newAdvertiserOnPost(ctx *fiber.Ctx) error {
	p := &newAdvertiserPresenter{}
	err := ctx.BodyParser(p)
	if err != nil {
		p.ErrorMessage = "One or more invalid inputs.Check and try again."
		return controllers.Render(ctx, adminNewAdvertiserTemplate, p)
	}
	email, err := user.NewEmail(p.Email)
	if err != nil {
		p.ErrorMessage = "Invalid email.Check and try again."
		return controllers.Render(ctx, adminNewAdvertiserTemplate, p)
	}
	ph, err := user.NewPhoneNumber(p.PhoneNumber)
	if err != nil {
		p.ErrorMessage = "Invalid phone number.Check and try again."
		return controllers.Render(ctx, adminNewAdvertiserTemplate, p)
	}
	err = c.advertiserInteractor.SignUpAdvertiser(email, ph, p.Name, p.Note)
	if err != nil {
		c.logger.Error("advertiser/signup", err)
		p.ErrorMessage = "Unable to sign up advertiser."
		return controllers.Render(ctx, adminNewAdvertiserTemplate, p)
	}
	return ctx.Redirect("/admin/advertisers")
}

type acceptDepositPresenter struct {
	Amount       float64 `form:"amount"`
	Start        string  `form:"start"`
	End          string  `form:"end"`
	ErrorMessage string
}

func (c *AdminAdvertisersController) acceptDepositOnPost(ctx *fiber.Ctx) error {
	advertiserID := ctx.Params("advertiserId", "")
	if advertiserID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	p := &acceptDepositPresenter{}
	err := ctx.BodyParser(p)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	s, err := time.Parse("2006-01-02", p.Start)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	e, err := time.Parse("2006-01-02", p.End)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	err = c.depositInteractor.AcceptDeposit(shared.ID(advertiserID), p.Amount, s, e)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	return ctx.Redirect("/admin/advertisers/" + advertiserID)
}

type newCampaignPresenter struct {
	Name  string `form:"name"`
	Start string `form:"start"`
	End   string `form:"end"`
}

func (c *AdminAdvertisersController) newCampaignOnPost(ctx *fiber.Ctx) error {
	advertiserID := ctx.Params("advertiserId", "")
	if advertiserID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	p := &newCampaignPresenter{}
	err := ctx.BodyParser(p)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	s, err := time.Parse("2006-01-02", p.Start)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	e, err := time.Parse("2006-01-02", p.End)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	err = c.campaignInteractor.StartCampaign(shared.ID(advertiserID), p.Name, s, e)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	return ctx.Redirect("/admin/advertisers/" + advertiserID)
}

func (c *AdminAdvertisersController) quitCampaignOnPost(ctx *fiber.Ctx) error {
	advertiserID := ctx.Params("advertiserId", "")
	if advertiserID == "" {
		return ctx.Redirect("/pages/error.html")
	}
	campaignID := ctx.Params("campaignId", "")
	if advertiserID == "" {
		return ctx.Redirect("/pages/error.html")
	}
	err := c.campaignInteractor.QuitCampaign(shared.ID(advertiserID), shared.ID(campaignID))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	return ctx.Redirect("/admin/advertisers/" + advertiserID)
}
