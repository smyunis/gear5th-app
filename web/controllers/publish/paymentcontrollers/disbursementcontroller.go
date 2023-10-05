package paymentcontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/disbursement"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var reqDisbursementTemplate *template.Template
var reqDisbursementResultTemplate *template.Template
var confirmDisbursementTemplate *template.Template

func init() {
	reqDisbursementTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/payments/request-disbursement.html"))
	reqDisbursementResultTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/payments/request-disbursement-result.html"))
	confirmDisbursementTemplate = template.Must(
		controllers.CardMainLayoutTemplate().ParseFiles(
			"web/views/publish/payments/confirm-disbursement.html"))

}

var availablePaymentMethods = []string{
	"Commercial Bank of Ethiopia",
	"Bank of Abyssinia",
	"Hibret Bank",
	"telebirr",
}

type reqDisbursementPresenter struct {
	AvailablePaymentMethods []string
	PaymetMethod            string `form:"payment-method"`
	Account                 string `form:"account"`
	Fullname                string `form:"fullname"`
	PhoneNumber             string `form:"phone-number"`
	ErrorMessage            string
}

type DisbursementController struct {
	authMiddleware         middlewares.JwtAuthenticationMiddleware
	earningInteractor      paymentinteractors.EarningInteractor
	disbursementInteractor paymentinteractors.DisbursementInteractor
	logger                 application.Logger
}

func NewDisbursementController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	earningInteractor paymentinteractors.EarningInteractor,
	disbursementInteractor paymentinteractors.DisbursementInteractor,
	logger application.Logger) DisbursementController {
	return DisbursementController{
		authMiddleware,
		earningInteractor,
		disbursementInteractor,
		logger,
	}
}

func (c *DisbursementController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/payments/disbursement/:disbursementId/confirm", c.confirmDisbursementOnGet)
	(*router).Add(fiber.MethodGet, "/payments/disbursement/:disbursementId/reject", c.rejectDisbursementOnGet)

	(*router).Use("/payments/request-disbursement", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/payments/request-disbursement", c.requestDisbursementOnGet)
	(*router).Add(fiber.MethodPost, "/payments/request-disbursement", c.requestDisbursementOnPost)
}

func (c *DisbursementController) requestDisbursementOnGet(ctx *fiber.Ctx) error {
	p := &reqDisbursementPresenter{
		AvailablePaymentMethods: availablePaymentMethods,
	}
	return controllers.Render(ctx, reqDisbursementTemplate, p)
}

func (c *DisbursementController) requestDisbursementOnPost(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	p := &reqDisbursementPresenter{
		AvailablePaymentMethods: availablePaymentMethods,
	}
	err = ctx.BodyParser(p)
	if err != nil {
		p.ErrorMessage = "One or more invalid input. Check and try again."
		return controllers.Render(ctx, reqDisbursementTemplate, p)
	}

	phone, err := user.NewPhoneNumber(p.PhoneNumber)
	if err != nil {
		p.ErrorMessage = "Invalid phone number. Check and try again."
		return controllers.Render(ctx, reqDisbursementTemplate, p)
	}

	profile := disbursement.PaymentProfile{
		PaymentMethod: p.PaymetMethod,
		Account:       p.Account,
		FullName:      p.Fullname,
		PhoneNumber:   phone,
	}

	err = c.disbursementInteractor.RequestDisbursement(publisherID, profile)
	if err != nil {
		p.ErrorMessage = "We're unable to request disbursements at the moment. Try again later."
		c.logger.Error("disbursement/request", err)
		return controllers.Render(ctx, reqDisbursementTemplate, p)
	}

	return controllers.Render(ctx, reqDisbursementResultTemplate, nil)

}

func (c *DisbursementController) confirmDisbursementOnGet(ctx *fiber.Ctx) error {
	disbursementID := ctx.Params("disbursementId", "")
	if disbursementID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	token := ctx.Query("token", "")
	if token == "" {
		return ctx.Redirect("/pages/error.html")
	}

	err := c.disbursementInteractor.ConfirmDisbursement(shared.ID(disbursementID), token)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	return controllers.Render(ctx, confirmDisbursementTemplate, true)
}

func (c *DisbursementController) rejectDisbursementOnGet(ctx *fiber.Ctx) error {
	disbursementID := ctx.Params("disbursementId", "")
	if disbursementID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	token := ctx.Query("token", "")
	if token == "" {
		return ctx.Redirect("/pages/error.html")
	}

	err := c.disbursementInteractor.RejectDisbursement(shared.ID(disbursementID), token)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	return controllers.Render(ctx, confirmDisbursementTemplate, false)
}
