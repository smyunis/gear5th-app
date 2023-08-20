package identitycontrollers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/web/controllers"
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
			Method: fiber.MethodGet,
			Path:   "/identity/managed/:userID/reset-password",
		},
		interactor,
	}
}

// http://localhost:5071/publish/identity/managed/01H8A9QYKP27CTS7D8E5F48H9M/reset-password?token=01H8AA6M7EBX3TJBG60T2W1DD5%2B989d017463acd69ccb60484860d1ab6a90d04c4f401d8b309f47964e7192346f

func (c ResetPasswordController) ResetPassword(ctx *fiber.Ctx) error {
	// resetPass := &resetPassword{}
	// err := ctx.BodyParser(resetPass)
	token := ctx.Query("token", "")
	userID := ctx.Params("userID")
	if userID == "" {
		return c.SendProblemDetails(ctx, fiber.StatusBadRequest, "", "")
	}

	if token == "" {
		return c.SendProblemDetails(ctx, fiber.StatusBadRequest,
			"Missing Reset Password Token",
			"valid reset password token must be provided")
	}

	email, err := user.NewEmail("mymail@gmail.com")
	if err != nil {

		return ctx.Redirect("https://bing.com")
	}

	//TODO get password from form
	err = c.interactor.ResetPassword(email, "gokuisking", token)
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
				"your email is not verified")
		}
		return c.SendProblemDetails(ctx, fiber.StatusInternalServerError, "", "")
	}

	ctx.SendStatus(fiber.StatusCreated)
	return ctx.Redirect("https://www.google.com")
}
