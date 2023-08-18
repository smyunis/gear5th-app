package identitycontrollers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type MangedUserSigninPresenter struct {
	HasValidationErrors       bool
	HasNonValidationError     bool
	NonValidationErrorMessage string
	Email                     string
	Password                  string
}

type ManagedUserController struct {
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewManagedUserController(interactor manageduserinteractors.ManagedUserInteractor) ManagedUserController {
	return ManagedUserController{
		interactor,
	}
}

func (c *ManagedUserController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/signin", c.Get)
	(*router).Add(fiber.MethodPost, "/identity/signin", c.Post)
}

func (c ManagedUserController) Get(ctx *fiber.Ctx) error {
	return ctx.Render("publish/identity/signin", nil, "publish/layouts/main")
}

func (c ManagedUserController) Post(ctx *fiber.Ctx) error {

	userEmail := ctx.FormValue("email", "")
	password := ctx.FormValue("password", "")

	if userEmail == "" || password == "" {
		return ctx.Render("publish/identity/signin", MangedUserSigninPresenter{
			HasValidationErrors: true,
		}, "publish/layouts/main")

	}

	email, err := user.NewEmail(userEmail)
	if err != nil {
		p := MangedUserSigninPresenter{
			HasValidationErrors: true,
			Email:               userEmail,
			Password:            password,
		}
		return ctx.Render("publish/identity/signin", p, "publish/layouts/main")
	}

	token, err := c.interactor.SignIn(email, password)
	if err != nil {
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			p := MangedUserSigninPresenter{
				HasNonValidationError:     true,
				NonValidationErrorMessage: "Your email has not been verified yet. Click on the verification link sent to your email",
				Email:                     userEmail,
				Password:                  password,
			}
			return ctx.Render("publish/identity/signin", p, "publish/layouts/main")

		}
		if errors.Is(err, manageduserinteractors.ErrAuthorization) {
			p := MangedUserSigninPresenter{
				HasValidationErrors: true,
				Email:               userEmail,
				Password:            password,
			}
			return ctx.Render("publish/identity/signin", p, "publish/layouts/main")
		}

		p := MangedUserSigninPresenter{
			HasNonValidationError:     true,
			NonValidationErrorMessage: "We're unable to sign you at at the moment. Please try agian later",
			Email:                     userEmail,
			Password:                  password,
		}
		return ctx.Render("publish/identity/signin", p, "publish/layouts/main")
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "gear5th-access-token",
		Value: token,
	})

	// return ctx.RedirectToRoute("publish/identity/signin", p, "publish/layouts/main")
	return ctx.Redirect("https://www.google.com")

}
