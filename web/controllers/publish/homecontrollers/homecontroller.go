package homecontrollers

import (
	"html/template"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/earning"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var homeTemplate *template.Template

func init() {
	homeTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().Funcs(template.FuncMap{
			"sub": func(a, b float64) string {
				diff := b - a
				return strconv.FormatFloat(diff, 'f', 2, 64)
			},
			"perdiff": func(a, b float64) string {
				diff := a - b
				percent := (diff / a) * 100
				if a == 0 {
					percent = 100
				}
				return strconv.FormatFloat(percent, 'f', 1, 64)
			},
			"isup": func(a, b float64) bool {
				return b >= a
			},
		}).ParseFiles(
			"web/views/publish/home/home.html"))
}

type earningCardPresenter struct {
	Prev      float64
	Cur       float64
	CurLabel  string
	PrevLabel string
}

type homePresenter struct {
	Day                       earningCardPresenter
	SevenDays                 earningCardPresenter
	Month                     earningCardPresenter
	CurrentBalance            float64
	BalanceTresholdPercentage float64
}

type HomeController struct {
	authMiddleware    middlewares.JwtAuthenticationMiddleware
	earningInteractor paymentinteractors.EarningInteractor
}

func NewHomeController(
	authMiddleware middlewares.JwtAuthenticationMiddleware,
	earningInteractor paymentinteractors.EarningInteractor) HomeController {
	return HomeController{
		authMiddleware,
		earningInteractor,
	}
}

func (c *HomeController) AddRoutes(router *fiber.Router) {
	(*router).Use("/home", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/home", c.onGet)
}

func (c *HomeController) onGet(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	yesterdayEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, 0, -1), today)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	dayBeforeYesterdayEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, 0, -3), today.AddDate(0, 0, -2))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	lastSevenDaysEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, 0, -7), today.AddDate(0, 0, -1))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	prevSevenDaysEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, 0, -14), today.AddDate(0, 0, -8))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	lastMonthEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, -1, 0), today.AddDate(0, 0, -1))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	prevMonthEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, -2, 0), today.AddDate(0, -1, 0))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	currentBalance, err := c.earningInteractor.CurrentBalance(publisherID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	balancePer := (currentBalance / earning.DisbursementRequestTreshold) * 100
	if balancePer > 100 {
		balancePer = 100
	}

	p := homePresenter{

		CurrentBalance:            currentBalance,
		BalanceTresholdPercentage: balancePer,

		Day: earningCardPresenter{
			Prev:      dayBeforeYesterdayEarning,
			Cur:       yesterdayEarning,
			CurLabel:  "Yesterday",
			PrevLabel: "previous day",
		},

		SevenDays: earningCardPresenter{
			Prev:      prevSevenDaysEarning,
			Cur:       lastSevenDaysEarning,
			CurLabel:  "Last seven days",
			PrevLabel: "previous seven days",
		},
		Month: earningCardPresenter{
			Prev:      prevMonthEarning,
			Cur:       lastMonthEarning,
			CurLabel:  "Last month",
			PrevLabel: "previous month",
		},
	}

	return controllers.Render(ctx, homeTemplate, p)
}
