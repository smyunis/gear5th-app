package identitycontrollers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/cmd/web/controllers"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type VerifyEmailController struct {
	controllers.Controller
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewVerifyEmailController(interactor manageduserinteractors.ManagedUserInteractor) VerifyEmailController {
	return VerifyEmailController{
		controllers.Controller{
			Method: fiber.MethodGet,
			Path:   "/identity/managed/:userId/verify-email",
		},
		interactor,
	}
}

func (c VerifyEmailController) VerifyEmail(ctx *fiber.Ctx) error {
	//TODO This endpoint should redirect to a page upon failure and success

	successPageURL := "https://www.google.com/"
	failPageURL := "https://duckduckgo.com/"

	userId := ctx.Params("userId")
	token := ctx.Query("token")
	if userId == "" || token == "" {
		return ctx.Redirect(failPageURL)
	}

	uID := shared.ID(userId)
	err := c.interactor.VerifyEmail(uID, token)

	if err != nil {
		return ctx.Redirect(failPageURL)
	}

	return ctx.Redirect(successPageURL)
}
