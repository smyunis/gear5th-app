package identitycontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/web/controllers"
)

var verifyEmailTemplate *template.Template

func init() {
	verifyEmailTemplate = template.Must(
		controllers.MainLayoutTemplate().ParseFiles(
			"web/views/publish/layouts/central-card.html",
			"web/views/publish/identity/managed/verify-email.html"))
}

type VerifyEmailController struct {
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewVerifyEmailController(interactor manageduserinteractors.ManagedUserInteractor) VerifyEmailController {
	return VerifyEmailController{
		interactor,
	}
}

func (c *VerifyEmailController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/managed/:userId/verify-email", c.onGet)
}

type VerifyEmailPresenter struct {
	IsSuccessful bool
}

func (c VerifyEmailController) onGet(ctx *fiber.Ctx) error {

	userId := ctx.Params("userId")
	token := ctx.Query("token")
	if userId == "" || token == "" {
		p := VerifyEmailPresenter{
			IsSuccessful: false,
		}
		return controllers.Render(ctx, verifyEmailTemplate, p)
	}

	uID := shared.ID(userId)
	err := c.interactor.VerifyEmail(uID, token)

	if err != nil {
		p := VerifyEmailPresenter{
			IsSuccessful: false,
		}
		return controllers.Render(ctx, verifyEmailTemplate, p)
	}
	p := VerifyEmailPresenter{
		IsSuccessful: true,
	}
	return controllers.Render(ctx, verifyEmailTemplate, p)

}
