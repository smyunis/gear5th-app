package adminadvertisers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/advertiserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var adminAdvertisersTemplate *template.Template

func init() {
	adminAdvertisersTemplate = template.Must(
		controllers.AdminMainLayoutTemplate().ParseFiles(
			"web/views/admin/advertisers/advertisers.html"))
}

type adminAdvertisersPresenter struct {
}

type AdminAdvertisersController struct {
	authMiddleware     middlewares.AdminAuthenticationMiddleware
	campaignInteractor advertiserinteractors.CampaignInteractor
	depositInteractor  paymentinteractors.DepositInteractor
	logger             application.Logger
}

func NewAdminAdvertisersController(
	authMiddleware middlewares.AdminAuthenticationMiddleware,
	campaignInteractor advertiserinteractors.CampaignInteractor,
	depositInteractor paymentinteractors.DepositInteractor,
	logger application.Logger) AdminAdvertisersController {
	return AdminAdvertisersController{
		authMiddleware,
		campaignInteractor,
		depositInteractor,
		logger,
	}
}

func (c *AdminAdvertisersController) AddRoutes(router *fiber.Router) {
	(*router).Use("/advertisers", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/advertisers", c.advertisersOnGet)
}

func (c *AdminAdvertisersController) advertisersOnGet(ctx *fiber.Ctx) error {
	return controllers.Render(ctx, adminAdvertisersTemplate, nil)
}
