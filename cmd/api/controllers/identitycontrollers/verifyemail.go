package identitycontrollers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/cmd/api/controllers"
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

	userId := ctx.Params("userId")
	token := ctx.Query("token")
	if userId == "" || token == "" {
		return c.SendProblemDetails(ctx, fiber.StatusBadRequest, "", "")
	}

	uID := shared.ID(userId)
	err := c.interactor.VerifyEmail(uID, token)

	if err != nil {
		if errors.Is(err, manageduserinteractors.ErrInvalidToken) {
			return c.SendProblemDetails(ctx, fiber.StatusNotFound,
				"Invalid Token",
				"token is either invalid or has expired")
		}
		if errors.Is(err, manageduserinteractors.ErrEntityNotFound) {
			return c.SendProblemDetails(ctx, fiber.StatusNotFound,
				"User Not Found",
				"user with provided ID is not signed up")
		}
		return c.SendProblemDetails(ctx, fiber.StatusInternalServerError, "", "")
	}

	ctx.SendStatus(fiber.StatusOK)
	return ctx.SendString("Email verified")
}
