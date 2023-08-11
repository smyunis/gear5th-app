package identitycontrollers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/cmd/api/controllers"
	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type emailStr struct {
	Email string `json:"email"`
}

type RequestPasswordResetController struct {
	controllers.Controller
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewRequestPasswordResetController(interactor manageduserinteractors.ManagedUserInteractor) RequestPasswordResetController {
	return RequestPasswordResetController{
		controllers.Controller{
			fiber.MethodPost,
			"/identity/managed/request-password-reset",
		},
		interactor,
	}
}

func (c RequestPasswordResetController) RequestPasswordReset(ctx *fiber.Ctx) error {
	var emailBody emailStr
	err := ctx.BodyParser(&emailBody)
	if err != nil {
		return c.SendProblemDetails(ctx, fiber.StatusBadRequest, "", "")
	}

	email, err := user.NewEmail(emailBody.Email)
	if err != nil {
		return c.SendProblemDetails(ctx, fiber.StatusBadRequest, "Invalid Email",
			fmt.Sprintf("%s is not a valid email", email.String()))
	}

	err = c.interactor.RequestResetPassword(email)
	if err != nil {

		if errors.Is(err, application.ErrEntityNotFound) {
			return c.SendProblemDetails(ctx, fiber.StatusNotFound,
				"Email Not Registered",
				"there is no user who signed up provided email")
		}
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			return c.SendProblemDetails(ctx, fiber.StatusPreconditionRequired,
				"Unverified Email",
				fmt.Sprintf("email %s is not verified", email.String()))
		}
		return c.SendProblemDetails(ctx, fiber.StatusInternalServerError, "", "")
	}
	ctx.SendStatus(fiber.StatusOK)
	return ctx.Send(nil)
}
