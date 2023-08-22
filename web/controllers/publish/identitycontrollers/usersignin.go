package identitycontrollers

import (
	"errors"
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

var signintemplate *template.Template

func init() {
	signintemplate = template.Must(
		controllers.MainLayoutTemplate().ParseFiles(
			"web/views/publish/layouts/central-card.html",
			"web/views/publish/identity/signin.html"))

}

const validationErrorMessage = "That email and password combination didn't work. Try again."

type UserSigninPresenter struct {
	Email        string `form:"email"`
	Password     string `form:"password"`
	StaySignedIn bool   `form:"stay-signed-in"`
	ErrorMessage string
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
	(*router).Add(fiber.MethodGet, "/identity/signin", c.onGet)
	(*router).Add(fiber.MethodPost, "/identity/signin", c.onPost)
}

func (*UserSignInController) onGet(ctx *fiber.Ctx) error {
	return controllers.Render(ctx, signintemplate, &UserSigninPresenter{})
	// return c.renderSignInPage(ctx, &UserSigninPresenter{})
}

func (c *UserSignInController) onPost(ctx *fiber.Ctx) error {
	p := &UserSigninPresenter{}
	err := ctx.BodyParser(p)

	if err != nil {
		p.ErrorMessage = validationErrorMessage
		return c.renderSignInPage(ctx, p)
	}

	email, err := user.NewEmail(p.Email)
	if err != nil {
		p.ErrorMessage = validationErrorMessage
		return c.renderSignInPage(ctx, p)
	}

	token, err := c.interactor.SignIn(email, p.Password)
	if err != nil {
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			p.ErrorMessage = "Your email has not been verified yet. Click on the verification link sent to your email."
			return c.renderSignInPage(ctx, p)
		}
		if errors.Is(err, manageduserinteractors.ErrAuthorization) {
			p.ErrorMessage = validationErrorMessage
			return c.renderSignInPage(ctx, p)
		}

		p.ErrorMessage = "We're unable to sign you in at the moment. Try agian later."
		return c.renderSignInPage(ctx, p)
	}

	if p.StaySignedIn {
		ctx.Cookie(&fiber.Cookie{
			Name:     AccessTokenCookieName(),
			Value:    token,
			Path:     "/publish",
			SameSite: "Lax",
			Expires:  time.Now().Add(720 * time.Hour), // 30 days
			// Secure: true, //TODO Must be set in production once TLS is setup
		})
	}

	return ctx.Redirect("/publish/home")
}

func (*UserSignInController) renderSignInPage(ctx *fiber.Ctx, p *UserSigninPresenter) error {
	return controllers.Render(ctx, signintemplate, p)
}

func AccessTokenCookieName() string {
	return "gear5th-access-token"
}
