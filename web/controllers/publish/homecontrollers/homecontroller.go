package homecontrollers

import (
	"html/template"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/earning"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
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
				diff := b - a
				percent := (diff / a) * 100
				if a == 0 {
					percent = 100
				}
				return strconv.FormatFloat(math.Abs(percent), 'f', 1, 64)
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

type countPresenter struct {
	Counts []int64
	Days   []int
	Label  string
}

type homePresenter struct {
	Day                       earningCardPresenter
	SevenDays                 earningCardPresenter
	Month                     earningCardPresenter
	CurrentBalance            float64
	BalanceTresholdPercentage float64

	ImpressionsCount countPresenter
	AdClicksCount    countPresenter
}

type HomeController struct {
	authMiddleware    middlewares.JwtAuthenticationMiddleware
	earningInteractor paymentinteractors.EarningInteractor
	adsInteractor     adsinteractors.AdsInteractor
}

func NewHomeController(
	authMiddleware middlewares.JwtAuthenticationMiddleware,
	earningInteractor paymentinteractors.EarningInteractor,
	adsInteractor adsinteractors.AdsInteractor) HomeController {
	return HomeController{
		authMiddleware,
		earningInteractor,
		adsInteractor,
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

	dayEarning, err := c.dayEarnings(publisherID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	sevenDayEarning, err := c.sevenDayEarning(publisherID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	monthEarning, err := c.monthEarning(publisherID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	currentBalance, err := c.earningInteractor.CurrentBalance(publisherID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	balancePer := earning.PercentOfDisbursementTreshold(currentBalance)

	ic, err := c.impressionsCount(publisherID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	ac, err := c.adclicksCount(publisherID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	p := homePresenter{
		CurrentBalance:            currentBalance,
		BalanceTresholdPercentage: balancePer,
		Day:                       dayEarning,
		SevenDays:                 sevenDayEarning,
		Month:                     monthEarning,
		ImpressionsCount:          ic,
		AdClicksCount:             ac,
	}

	return controllers.Render(ctx, homeTemplate, p)
}

func (c *HomeController) monthEarning(publisherID shared.ID) (earningCardPresenter, error) {

	today,_ := shared.DayTimeEdges(time.Now())

	lastMonthEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, -1, 0), today)
	if err != nil {
		return earningCardPresenter{}, err
	}

	prevMonthEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, -2, 0), today.AddDate(0, -1, 0))
	if err != nil {
		return earningCardPresenter{}, err
	}
	return earningCardPresenter{
		Prev:      prevMonthEarning,
		Cur:       lastMonthEarning,
		CurLabel:  "Last month",
		PrevLabel: "previous month",
	}, nil
}

func (c *HomeController) sevenDayEarning(publisherID shared.ID) (earningCardPresenter, error) {

	today,_ := shared.DayTimeEdges(time.Now())

	lastSevenDaysEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, 0, -7), today)
	if err != nil {
		return earningCardPresenter{}, err
	}

	prevSevenDaysEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, 0, -14), today.AddDate(0, 0, -7))
	if err != nil {
		return earningCardPresenter{}, err
	}
	return earningCardPresenter{
		Prev:      prevSevenDaysEarning,
		Cur:       lastSevenDaysEarning,
		CurLabel:  "Last seven days",
		PrevLabel: "previous seven days",
	}, nil
}

func (c *HomeController) dayEarnings(publisherID shared.ID) (earningCardPresenter, error) {

	today,_ := shared.DayTimeEdges(time.Now())

	yesterdayEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, 0, -1), today)
	if err != nil {
		return earningCardPresenter{}, err
	}

	dayBeforeYesterdayEarning, err := c.earningInteractor.Earnings(publisherID,
		today.AddDate(0, 0, -2), today.AddDate(0, 0, -1))
	if err != nil {
		return earningCardPresenter{}, err
	}
	return earningCardPresenter{
		Prev:      dayBeforeYesterdayEarning,
		Cur:       yesterdayEarning,
		CurLabel:  "Yesterday",
		PrevLabel: "previous day",
	}, nil
}

func (c *HomeController) impressionsCount(publisherID shared.ID) (countPresenter, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var impCount []int64 = make([]int64, 30)
	var days []int = make([]int, 30)
	var labelMonths []string

	for i := 0; i > -30; i-- {
		count, err := c.adsInteractor.ImpressionsCount(publisherID,
			today.AddDate(0, 0, i-1), today.AddDate(0, 0, i))
		if err != nil {
			continue
		}
		impCount[i*-1] = count
		days[i*-1] = today.AddDate(0, 0, i-1).Day()

		m := today.AddDate(0, 0, i-1).Format("Jan")
		if !slices.Contains(labelMonths, m) {
			labelMonths = append(labelMonths, m)
		}
	}

	slices.Reverse(impCount)
	slices.Reverse(days)
	slices.Reverse(labelMonths)

	return countPresenter{
		Counts: impCount,
		Days:   days,
		Label:  strings.Join(labelMonths, "/"),
	}, nil
}

func (c *HomeController) adclicksCount(publisherID shared.ID) (countPresenter, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	var clicksCount []int64 = make([]int64, 30)
	var days []int = make([]int, 30)
	var labelMonths []string

	for i := 0; i > -30; i-- {
		count, err := c.adsInteractor.AdClicksCount(publisherID,
			today.AddDate(0, 0, i-1), today.AddDate(0, 0, i))
		if err != nil {
			continue
		}
		clicksCount[i*-1] = count
		days[i*-1] = today.AddDate(0, 0, i-1).Day()

		m := today.AddDate(0, 0, i-1).Format("Jan")
		if !slices.Contains(labelMonths, m) {
			labelMonths = append(labelMonths, m)
		}
	}

	slices.Reverse(clicksCount)
	slices.Reverse(days)
	slices.Reverse(labelMonths)

	return countPresenter{
		Counts: clicksCount,
		Days:   days,
		Label:  strings.Join(labelMonths, "/"),
	}, nil
}
