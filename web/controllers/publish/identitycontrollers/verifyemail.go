package identitycontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

var verifyEmailTemplate *template.Template

func init() {
	verifyEmailTemplate = template.Must(
		controllers.CardMainLayoutTemplate().ParseFiles(
			"web/views/publish/identity/managed/verify-email.html"))
}

type verifyEmailPresenter struct {
	IsSuccessful bool
}

type VerifyEmailController struct {
	interactor identityinteractors.ManagedUserInteractor

	logger application.Logger
}

func NewVerifyEmailController(interactor identityinteractors.ManagedUserInteractor,
	logger application.Logger) VerifyEmailController {
	return VerifyEmailController{
		interactor,
		logger,
	}
}

func (c *VerifyEmailController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/managed/:userId/verify-email", c.onGet)
}

func (c *VerifyEmailController) onGet(ctx *fiber.Ctx) error {

	userId := ctx.Params("userId", "")
	token := ctx.Query("token", "")
	if userId == "" || token == "" {
		p := verifyEmailPresenter{
			IsSuccessful: false,
		}
		return controllers.Render(ctx, verifyEmailTemplate, p)
	}

	uID := shared.ID(userId)
	err := c.interactor.VerifyEmail(uID, token)

	if err != nil {
		// c.logger.Error("identity/verifyemail", err)
		p := verifyEmailPresenter{
			IsSuccessful: false,
		}
		return controllers.Render(ctx, verifyEmailTemplate, p)
	}
	p := verifyEmailPresenter{
		IsSuccessful: true,
	}
	return controllers.Render(ctx, verifyEmailTemplate, p)

}
