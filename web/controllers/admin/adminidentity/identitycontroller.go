package adminidentity

import (
	"errors"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

var adminSignInTemplate *template.Template

const validationErrorMessage = "That email and password combination didn't work. Try again."

func init() {
	adminSignInTemplate = template.Must(
		controllers.AdminCardLayoutTemplate().ParseFiles(
			"web/views/admin/identity/signin.html"))
}

type adminSignInPresenter struct {
	Email        string `form:"email"`
	Password     string `form:"password"`
	ErrorMessage string
}

type AdminIdentityController struct {
	managedUserInteractor identityinteractors.ManagedUserInteractor
	logger                application.Logger
}

func NewAdminIdentityController(
	interactor identityinteractors.ManagedUserInteractor,
	logger application.Logger) AdminIdentityController {
	return AdminIdentityController{
		interactor,
		logger,
	}
}

func (c *AdminIdentityController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/signin", c.signInOnGet)
	(*router).Add(fiber.MethodPost, "/identity/signin", c.signInOnPost)
}

func (c *AdminIdentityController) signInOnGet(ctx *fiber.Ctx) error {
	return controllers.Render(ctx, adminSignInTemplate, nil)
}

func (c *AdminIdentityController) signInOnPost(ctx *fiber.Ctx) error {
	p := &adminSignInPresenter{}

	err := ctx.BodyParser(p)

	if err != nil {
		p.ErrorMessage = "One or more invalid inputs. Check and try again."
		return controllers.Render(ctx, adminSignInTemplate, p)
	}

	email, err := user.NewEmail(p.Email)
	if err != nil {
		p.ErrorMessage = validationErrorMessage
		return controllers.Render(ctx, adminSignInTemplate, p)
	}

	token, err := c.managedUserInteractor.AdminSignIn(email, p.Password)
	if err != nil {
		if errors.Is(err, application.ErrAuthorization) {
			p.ErrorMessage = validationErrorMessage
			return controllers.Render(ctx, adminSignInTemplate, p)
		}

		c.logger.Error("admin/identity/signin", err)
		p.ErrorMessage = "We're unable to sign you in at the moment. Try again later."
		return controllers.Render(ctx, adminSignInTemplate, p)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     controllers.AdminAccessTokenCookieName,
		Value:    token,
		Path:     "/admin",
		SameSite: "Lax",
		// Secure: true, //TODO Must be set in production once TLS is setup
	})

	return ctx.Redirect("/admin/dashboard")
}
