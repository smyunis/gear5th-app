package identitycontrollers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type resetPasswordPresenter struct {
	Email        string `form:"email"`
	NewPassword  string `form:"new-password"`
	Token        string `form:"token"`
	UserID       string `form:"user-id"`
	ErrorMessage string
}

type ResetPasswordController struct {
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewResetPasswordController(interactor manageduserinteractors.ManagedUserInteractor) ResetPasswordController {
	return ResetPasswordController{
		interactor,
	}
}

func (c *ResetPasswordController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/managed/:userID/reset-password", c.onGet)
	(*router).Add(fiber.MethodPost, "/identity/managed/:userID/reset-password", c.onPost)
}

func (c ResetPasswordController) onGet(ctx *fiber.Ctx) error {

	token := ctx.Query("token", "")
	userID := ctx.Params("userID")
	if userID == "" || token == "" {
		return ctx.Render("publish/identity/managed/reset-password-result", false, "publish/layouts/main")
	}

	presenter := &resetPasswordPresenter{
		Token: token,
		UserID: userID,
	}
	return ctx.Render("publish/identity/managed/reset-password", presenter, "publish/layouts/main")

}

func (c ResetPasswordController) onPost(ctx *fiber.Ctx) error {

	presenter := &resetPasswordPresenter{}
	err := ctx.BodyParser(presenter)
	if err != nil {
		presenter.ErrorMessage = "There are one or more invalid inputs. Check and try again"
		return ctx.Render("publish/identity/managed/reset-password", presenter, "publish/layouts/main")

	}

	email, err := user.NewEmail(presenter.Email)
	if err != nil {
		presenter.ErrorMessage = presenter.Email + " is not a valid email. Check and try again."
		return ctx.Render("publish/identity/managed/reset-password", presenter, "publish/layouts/main")
	}

	err = c.interactor.ResetPassword(email, presenter.NewPassword, presenter.Token)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			presenter.ErrorMessage = "There is no user who signed up with that email. Check and try agian."
			return ctx.Render("publish/identity/managed/reset-password", presenter, "publish/layouts/main")
		}
		if errors.Is(err, manageduserinteractors.ErrInvalidToken) {
			presenter.ErrorMessage = "We're unable to reset your password. This may be due to the link sent to your email being altered or you have entered a wrong email. Check and try again."
			return ctx.Render("publish/identity/managed/reset-password", presenter, "publish/layouts/main")
		}
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			presenter.ErrorMessage = "Your email has not been verified by our system. Click on a verification link sent to your email then try resetting your password again."
			return ctx.Render("publish/identity/managed/reset-password", presenter, "publish/layouts/main")
		}
		presenter.ErrorMessage = "We're unable to reset your password at the moment. Try again later."
		return ctx.Render("publish/identity/managed/reset-password", presenter, "publish/layouts/main")
	}

	return ctx.Render("publish/identity/managed/reset-password-result", true, "publish/layouts/main")
}
