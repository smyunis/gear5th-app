package paymentcontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var reqDisbursementTemplate *template.Template

func init() {
	reqDisbursementTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/payments/request-disbursement.html"))

}

type reqDisbursementPresenter struct {
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
	(*router).Use("/payments/request-disbursement", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/payments/request-disbursement", c.requestDisbursementOnGet)
}

func (c *DisbursementController) requestDisbursementOnGet(ctx *fiber.Ctx) error {
	return controllers.Render(ctx, reqDisbursementTemplate, nil)
}
