package adminpublisherpayments

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/paymentinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/disbursement"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var adminPublisherPaymentsTemplate *template.Template
var adminPubDisbursementTemplate *template.Template

// var adminSettleDisbursementTemplate *template.Template

func init() {
	adminPublisherPaymentsTemplate = template.Must(
		controllers.AdminMainLayoutTemplate().ParseFiles(
			"web/views/admin/publishers/payments.html"))
	adminPubDisbursementTemplate = template.Must(
		controllers.AdminMainLayoutTemplate().ParseFiles(
			"web/views/admin/publishers/disbursement.html"))
}

type adminPublisherPaymentsPresenter struct {
	ConfirmedDisbursements []disbursement.Disbursement
}

type AdminPublisherPaymentsController struct {
	authMiddleware         middlewares.AdminAuthenticationMiddleware
	adsInteractor          adsinteractors.AdsInteractor
	depositInteractor      paymentinteractors.DepositInteractor
	disbursementInteractor paymentinteractors.DisbursementInteractor
	earningInteractor      paymentinteractors.EarningInteractor
	logger                 application.Logger
}

func NewAdminPublisherPaymentsController(
	authMiddleware middlewares.AdminAuthenticationMiddleware,
	adsInteractor adsinteractors.AdsInteractor,
	depositInteractor paymentinteractors.DepositInteractor,
	disbursementInteractor paymentinteractors.DisbursementInteractor,
	earningInteractor paymentinteractors.EarningInteractor,
	logger application.Logger) AdminPublisherPaymentsController {
	return AdminPublisherPaymentsController{
		authMiddleware,
		adsInteractor,
		depositInteractor,
		disbursementInteractor,
		earningInteractor,
		logger,
	}
}

func (c *AdminPublisherPaymentsController) AddRoutes(router *fiber.Router) {
	(*router).Use("/publishers/payments", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/publishers/payments", c.pubPaymentOnGet)
	(*router).Add(fiber.MethodGet, "/publishers/payments/disbursement/:disbursementId", c.disbursementOnGet)
	(*router).Add(fiber.MethodPost, "/publishers/payments/disbursement/:disbursementId/settle", c.settleOnPost)
}

func (c *AdminPublisherPaymentsController) pubPaymentOnGet(ctx *fiber.Ctx) error {

	p := &adminPublisherPaymentsPresenter{}
	confirmedDisbursements, err := c.disbursementInteractor.DisbursementsWithStatus(disbursement.Confirmed)
	if err != nil {
		p.ConfirmedDisbursements = make([]disbursement.Disbursement, 0)
	}

	p.ConfirmedDisbursements = confirmedDisbursements

	return controllers.Render(ctx, adminPublisherPaymentsTemplate, p)
}

func (c *AdminPublisherPaymentsController) disbursementOnGet(ctx *fiber.Ctx) error {
	disbursementID := ctx.Params("disbursementId", "")
	if disbursementID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	d, err := c.disbursementInteractor.Disbursement(shared.ID(disbursementID))
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	return controllers.Render(ctx, adminPubDisbursementTemplate, d)
}


func (c *AdminPublisherPaymentsController) settleOnPost(ctx *fiber.Ctx) error {
	disbursementID := ctx.Params("disbursementId", "")
	if disbursementID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	disbursementForm := ctx.FormValue("settle", "")
	if disbursementForm != disbursementID {
		return ctx.Redirect("/pages/error.html")
	}
	settlementRemark := ctx.FormValue("set-remark", "")
	if settlementRemark == "" {
		return ctx.Redirect("/pages/error.html")
	}

	err := c.disbursementInteractor.SettleDisbursement(shared.ID(disbursementID), settlementRemark)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	return ctx.Redirect("/admin/publishers/payments")
}
