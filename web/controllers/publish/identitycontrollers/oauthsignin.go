package identitycontrollers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

type oAuthSignInPresenter struct {
	ClientID     string `form:"clientId"`
	Credential   string `form:"credential"`
	SelectBy     string `form:"select_by"`
	CSRFToken    string `form:"g_csrf_token"`
	ErrorMessage string
	Email        string
	Password     string
}

type OAuthSignInController struct {
	interactor identityinteractors.OAuthUserInteractor
	logger     application.Logger
}

func NewOAuthSignInController(interactor identityinteractors.OAuthUserInteractor,
	logger application.Logger) OAuthSignInController {
	return OAuthSignInController{
		interactor,
		logger,
	}
}

func (c *OAuthSignInController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodPost, "/identity/oauth/signin", c.onPost)
}

func (c *OAuthSignInController) onPost(ctx *fiber.Ctx) error {
	p := &oAuthSignInPresenter{}
	err := ctx.BodyParser(p)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	googleCSRFCookie := ctx.Cookies("g_csrf_token", "")
	if p.CSRFToken != googleCSRFCookie {
		c.logger.Info("oauth/signin", "CSRF mismatch")
		return ctx.Redirect("/pages/error.html")
	}

	token, err := c.interactor.SignIn(p.Credential)
	if err != nil {
		if errors.Is(err, identityinteractors.ErrAuthorization) {
			p.ErrorMessage = "We're unable to sign you in. Make sure you have signed up before attempting to sign in."
			return controllers.Render(ctx, signintemplate, p)
		}

		c.logger.Error("identity/oauth/signin", err)
		p.ErrorMessage = "We're unable to sign you in at the moment. Try again later."
		return controllers.Render(ctx, signintemplate, p)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     controllers.AccessTokenCookieName,
		Value:    token,
		Path:     "/publish",
		SameSite: "Lax",
		Expires:  time.Now().Add(720 * time.Hour), // 30 days
		// Secure: true, //TODO Must be set in production once TLS is setup
	})

	return ctx.Redirect("/publish/home")
}
