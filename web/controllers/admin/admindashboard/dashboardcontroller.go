package admindashboard

import (
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/earning"
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
	YesterdayRPI             float64
	RPIInfation              bool
	ProfitMargin             float64
	Profit                   float64
}

type AdminDashboardController struct {
	authMiddleware    middlewares.AdminAuthenticationMiddleware
	adsInteractor     adsinteractors.AdsInteractor
	depositInteractor paymentinteractors.DepositInteractor
	earningInteractor paymentinteractors.EarningInteractor
	logger            application.Logger
}

func NewAdminDashboardController(
	authMiddleware middlewares.AdminAuthenticationMiddleware,
	adsInteractor adsinteractors.AdsInteractor,
	depositInteractor paymentinteractors.DepositInteractor,
	earningInteractor paymentinteractors.EarningInteractor,
	logger application.Logger) AdminDashboardController {
	return AdminDashboardController{
		authMiddleware,
		adsInteractor,
		depositInteractor,
		earningInteractor,
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

	rpi, err := c.earningInteractor.DailyRatePerImpression(yesterday)
	profit := fund - (rpi * float64(impCount))
	profMargin := (profit / fund) * 100
	if fund == 0 {
		profMargin = 0.0
	}

	p := &adminDashboardPresenter{
		YesterdayImpressionCount: impCount,
		YesterdayDailyFund:       fund,
		YesterdayRPI:             rpi,
		RPIInfation:              rpi < earning.FixedRatePerImpression,
		ProfitMargin:             profMargin,
		Profit:                   profit,
	}

	return controllers.Render(ctx, adminDashboardTemplate, p)
}
