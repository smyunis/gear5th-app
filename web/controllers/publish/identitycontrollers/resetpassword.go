package identitycontrollers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/web/controllers"
	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type resetPassword struct {
	Token       string `json:"token"`
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
}

type ResetPasswordController struct {
	controllers.Controller
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewResetPasswordController(interactor manageduserinteractors.ManagedUserInteractor) ResetPasswordController {
	return ResetPasswordController{
		controllers.Controller{
			Method: fiber.MethodPost,
			Path:   "/identity/managed/:userID/reset-password",
		},
		interactor,
	}
}

func (c ResetPasswordController) ResetPassword(ctx *fiber.Ctx) error {
	resetPass := &resetPassword{}
	err := ctx.BodyParser(resetPass)
	userID := ctx.Params("userID")
	if err != nil || userID == "" {
		return c.SendProblemDetails(ctx, fiber.StatusBadRequest, "", "")
	}

	if resetPass.Token == "" {
		return c.SendProblemDetails(ctx, fiber.StatusBadRequest,
			"Missing Reset Password Token",
			"valid reset password token must be provided")
	}

	email, err := user.NewEmail(resetPass.Email)
	if err != nil {
		return c.SendProblemDetails(ctx, fiber.StatusBadRequest,
			"Invalid Email",
			fmt.Sprintf("%s is not a valid email", email.String()))
	}

	err = c.interactor.ResetPassword(email, resetPass.NewPassword, resetPass.Token)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			return c.SendProblemDetails(ctx, fiber.StatusNotFound,
				"Email Not Registered",
				"there is no user who signed up provided email")
		}
		if errors.Is(err, manageduserinteractors.ErrInvalidToken) {
			return c.SendProblemDetails(ctx, fiber.StatusBadRequest,
				"Invalid Reset Password Token",
				"the reset password token provided is invalid or has expired")
		}
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			return c.SendProblemDetails(ctx, fiber.StatusPreconditionRequired,
				"Unverified Email",
				fmt.Sprintf("email %s is not verified", email.String()))
		}
		return c.SendProblemDetails(ctx, fiber.StatusInternalServerError, "", "")
	}

	ctx.SendStatus(fiber.StatusCreated)
	return ctx.Send(nil)
}
