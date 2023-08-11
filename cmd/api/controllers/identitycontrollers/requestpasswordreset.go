package identitycontrollers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/cmd/api/controllers"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/pkg/problemdetails"
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
		ctx.SendStatus(fiber.StatusBadRequest)
		return ctx.JSON(problemdetails.NewProblemDetails(fiber.StatusBadRequest))
	}

	email, err := user.NewEmail(emailBody.Email)
	if err != nil {
		probDetails := problemdetails.NewProblemDetails(fiber.StatusBadRequest)
		probDetails.Title = "Invalid Email"
		probDetails.Detail = fmt.Sprintf("%s is not a valid email", email.String())
		ctx.SendStatus(fiber.StatusBadRequest)
		return ctx.JSON(probDetails)
	}

	err = c.interactor.RequestResetPassword(email)
	if err != nil {
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			prob := problemdetails.NewProblemDetails(fiber.StatusPreconditionRequired)
			prob.Title = "Unverified Email"
			prob.Detail = fmt.Sprintf("email %s is not verified", email.String())
			ctx.SendStatus(fiber.StatusPreconditionRequired)
			return ctx.JSON(prob)
		}
		ctx.SendStatus(fiber.StatusInternalServerError)
		return ctx.JSON(problemdetails.NewProblemDetails(fiber.StatusInternalServerError))
	}
	ctx.SendStatus(fiber.StatusOK)
	return ctx.Send(nil)
}
