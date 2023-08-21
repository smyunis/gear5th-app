package identitycontrollers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type VerifyEmailController struct {
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewVerifyEmailController(interactor manageduserinteractors.ManagedUserInteractor) VerifyEmailController {
	return VerifyEmailController{
		interactor,
	}
}

func (c *VerifyEmailController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/managed/:userId/verify-email", c.VerifyEmail)
}

type VerifyEmailPresenter struct {
	IsSuccessful bool
}
// http://localhost:5071/publish/identity/managed/01H8C1BQ67BRSKPH6Q1WE3T856/verify-email?token=01H8C1BQETRNASJX65Q7B8W7Y6%2Be1126967cddb1afd1f157d16ae43ba893eb28c665d553b2442652080fe7f70ff
func (c VerifyEmailController) VerifyEmail(ctx *fiber.Ctx) error {

	userId := ctx.Params("userId")
	token := ctx.Query("token")
	if userId == "" || token == "" {
		p := VerifyEmailPresenter{
			IsSuccessful: false,
		}
		return ctx.Render("publish/identity/managed/verify-email", p, "publish/layouts/main")
	}

	uID := shared.ID(userId)
	err := c.interactor.VerifyEmail(uID, token)

	if err != nil {
		p := VerifyEmailPresenter{
			IsSuccessful: false,
		}
		return ctx.Render("publish/identity/managed/verify-email", p, "publish/layouts/main")
	}
	p := VerifyEmailPresenter{
		IsSuccessful: true,
	}
	return ctx.Render("publish/identity/managed/verify-email", p, "publish/layouts/main")
}
