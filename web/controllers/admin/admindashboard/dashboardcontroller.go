package admindashboard

import (
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var adminDashboardTemplate *template.Template

func init() {
	adminDashboardTemplate = template.Must(
		controllers.AdminMainLayoutTemplate().ParseFiles(
			"web/views/admin/dashboard/dashboard.html"))
}

type adminDashboardPresenter struct {
	YesterdayImpressionCount int
	YesterdayDailyFund       float64
}

type AdminDashboardController struct {
	authMiddleware    middlewares.AdminAuthenticationMiddleware
	adsInteractor     adsinteractors.AdsInteractor
	depositInteractor paymentinteractors.DepositInteractor
	logger            application.Logger
}

func NewAdminDashboardController(
	authMiddleware middlewares.AdminAuthenticationMiddleware,
	adsInteractor adsinteractors.AdsInteractor,
	depositInteractor paymentinteractors.DepositInteractor,
	logger application.Logger) AdminDashboardController {
	return AdminDashboardController{
		authMiddleware,
		adsInteractor,
		depositInteractor,
		logger,
	}
}

func (c *AdminDashboardController) AddRoutes(router *fiber.Router) {
	(*router).Use("/dashboard", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/dashboard", c.dashboardOnGet)
}

func (c *AdminDashboardController) dashboardOnGet(ctx *fiber.Ctx) error {

	yesterday := time.Now().AddDate(0, 0, -1)

	impCount, err := c.adsInteractor.TotalImpressionCount(yesterday)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	fund, err := c.depositInteractor.TotalDailyFund(yesterday)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	p := &adminDashboardPresenter{
		YesterdayImpressionCount: impCount,
		YesterdayDailyFund:       fund,
	}

	return controllers.Render(ctx, adminDashboardTemplate, p)
}
