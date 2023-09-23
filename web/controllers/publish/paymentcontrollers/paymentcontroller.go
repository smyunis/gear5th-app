package paymentcontrollers

import (
	"html/template"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/disbursement"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/earning"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var paymentTemplate *template.Template

func init() {
	paymentTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/payments/payments.html"))

}

type paymentPresenter struct {
	CurrentBalance                   string
	DisbursementTreshold             string
	PercentageOfDisbursementTreshold string
	CanRequestDisbursement            bool
	SettledDisbursements             []disbursement.Disbursement
}

type PaymentController struct {
	authMiddleware         middlewares.JwtAuthenticationMiddleware
	earningInteractor      paymentinteractors.EarningInteractor
	disbursementInteractor paymentinteractors.DisbursementInteractor
	logger                 application.Logger
}

func NewPaymentController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	earningInteractor paymentinteractors.EarningInteractor,
	disbursementInteractor paymentinteractors.DisbursementInteractor,
	logger application.Logger) PaymentController {
	return PaymentController{
		authMiddleware,
		earningInteractor,
		disbursementInteractor,
		logger,
	}
}

func (c *PaymentController) AddRoutes(router *fiber.Router) {
	(*router).Use("/payments", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/payments", c.paymentOnGet)
}

func (c *PaymentController) paymentOnGet(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	balance, err := c.earningInteractor.CurrentBalance(publisherID)
	if err != nil {
		c.logger.Error("earning/balance", err)
		return ctx.Redirect("/pages/error.html")
	}

	p := &paymentPresenter{}

	prevDisbursements, err := c.disbursementInteractor.SettledDisbursements(publisherID)
	if err != nil {
		p.SettledDisbursements = make([]disbursement.Disbursement, 0)
	}

	p.SettledDisbursements = prevDisbursements
	p.CurrentBalance = strconv.FormatFloat(balance, 'f', 2, 64)
	p.DisbursementTreshold = strconv.Itoa(earning.DisbursementRequestTreshold)
	percentage := (balance / float64(earning.DisbursementRequestTreshold)) * 100
	if percentage >= 100 {
		percentage = 100
	}
	p.PercentageOfDisbursementTreshold = strconv.FormatFloat(percentage, 'f', 2, 64)
	p.CanRequestDisbursement = c.earningInteractor.CanRequestDisbursement(publisherID)

	return controllers.Render(ctx, paymentTemplate, p)
}
