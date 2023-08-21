package identitycontrollers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type requestPasswordResetPresenter struct {
	Email string
	ErrorMessage string
}

type RequestPasswordResetController struct {
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewRequestPasswordResetController(interactor manageduserinteractors.ManagedUserInteractor) RequestPasswordResetController {
	return RequestPasswordResetController{
		interactor,
	}
}

func (c *RequestPasswordResetController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/managed/request-password-reset", c.onGet)
	(*router).Add(fiber.MethodPost, "/identity/managed/request-password-reset", c.onPost)
}

func (c RequestPasswordResetController) onGet(ctx *fiber.Ctx) error {
	return ctx.Render("publish/identity/managed/request-password-reset", nil, "publish/layouts/main")
}

func (c RequestPasswordResetController) onPost(ctx *fiber.Ctx) error {
	e := ctx.FormValue("email", "")
	if e == "" {
		p := requestPasswordResetPresenter{
			Email: e,
			ErrorMessage: "Invalid email. Check and try again.",
		}
		return ctx.Render("publish/identity/managed/request-password-reset", p, "publish/layouts/main")
	}

	email, err := user.NewEmail(e)
	if err != nil {
		p := requestPasswordResetPresenter{
			ErrorMessage: fmt.Sprintf("%s is not a valid email. Check and try agian.", e),
		}
		return ctx.Render("publish/identity/managed/request-password-reset", p, "publish/layouts/main")
	}

	err = c.interactor.RequestResetPassword(email)
	if err != nil {

		if errors.Is(err, application.ErrEntityNotFound) {
			p := requestPasswordResetPresenter{
				Email: e,
				ErrorMessage: "There is no user who signed up with that email. Check and try agian.",
			}
			return ctx.Render("publish/identity/managed/request-password-reset", p, "publish/layouts/main")
		}
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			p := requestPasswordResetPresenter{
				Email: e,
				ErrorMessage: "Your email has not been verified by our system. Click on a verification link sent to your email then try resetting your password again.",
			}
			return ctx.Render("publish/identity/managed/request-password-reset", p, "publish/layouts/main")
		}
		p := requestPasswordResetPresenter{
			Email: e,
			ErrorMessage: "We're unable to reset your password at the moment. Try again later.",
		}
		return ctx.Render("publish/identity/managed/request-password-reset", p, "publish/layouts/main")
	}
	p := requestPasswordResetPresenter{
		Email: e,
	}
	return ctx.Render("publish/identity/managed/request-password-reset-success", p, "publish/layouts/main")
}
