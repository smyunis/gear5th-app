package identitycontrollers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type UserSigninPresenter struct {
	HasValidationError        bool
	HasNonValidationError     bool
	NonValidationErrorMessage string
	Email                     string
	Password                  string
}

type UserSignInController struct {
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewUserSignInController(interactor manageduserinteractors.ManagedUserInteractor) UserSignInController {
	return UserSignInController{
		interactor,
	}
}

func (c *UserSignInController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/signin", c.userSignInOnGet)
	(*router).Add(fiber.MethodPost, "/identity/signin", c.userSignInOnPost)
}

func (*UserSignInController) userSignInOnGet(ctx *fiber.Ctx) error {
	return ctx.Render("publish/identity/signin", nil, "publish/layouts/main")
}

func (c *UserSignInController) userSignInOnPost(ctx *fiber.Ctx) error {

	userEmail := ctx.FormValue("email", "")
	password := ctx.FormValue("password", "")
	staySignedIn := c.shouldStaySignedIn(ctx.FormValue("stay-signed-in", ""))

	if userEmail == "" || password == "" {
		return ctx.Render("publish/identity/signin", UserSigninPresenter{
			HasValidationError: true,
		}, "publish/layouts/main")

	}

	email, err := user.NewEmail(userEmail)
	if err != nil {
		p := UserSigninPresenter{
			HasValidationError: true,
			Email:              userEmail,
			Password:           password,
		}
		return ctx.Render("publish/identity/signin", p, "publish/layouts/main")
	}

	token, err := c.interactor.SignIn(email, password)
	if err != nil {
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			p := UserSigninPresenter{
				HasNonValidationError:     true,
				NonValidationErrorMessage: "Your email has not been verified yet. Click on the verification link sent to your email.",
				Email:                     userEmail,
				Password:                  password,
			}
			return ctx.Render("publish/identity/signin", p, "publish/layouts/main")

		}
		if errors.Is(err, manageduserinteractors.ErrAuthorization) {
			p := UserSigninPresenter{
				HasValidationError: true,
				Email:              userEmail,
				Password:           password,
			}
			return ctx.Render("publish/identity/signin", p, "publish/layouts/main")
		}

		p := UserSigninPresenter{
			HasNonValidationError:     true,
			NonValidationErrorMessage: "We're unable to sign you in at the moment. Try agian later.",
			Email:                     userEmail,
			Password:                  password,
		}
		return ctx.Render("publish/identity/signin", p, "publish/layouts/main")
	}

	if staySignedIn {
		ctx.Cookie(&fiber.Cookie{
			Name:     AccessTokenCookieName(),
			Value:    token,
			Path:     "/publish",
			SameSite: "Lax",
			Expires:  time.Now().Add(720 * time.Hour), // 30 days
			// Secure: true, //TODO Must be set in production once TLS is setup
		})
	}

	return ctx.Redirect("https://www.google.com")
	// return ctx.RedirectToRoute("publish/identity/signin", p, "publish/layouts/main")

}

func (*UserSignInController) shouldStaySignedIn(formData string) bool {
	return formData == "on"
}

func AccessTokenCookieName() string {
	return "gear5th-access-token"
}
