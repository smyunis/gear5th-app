package identitycontrollers

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

var requestPasswordResetTemplate *template.Template
var requestPasswordResetSuccessTemplate *template.Template

func init() {
	requestPasswordResetTemplate = template.Must(
		controllers.MainLayoutTemplate().ParseFiles(
			"web/views/publish/layouts/central-card.html",
			"web/views/publish/identity/managed/request-password-reset.html"))

	requestPasswordResetSuccessTemplate = template.Must(
		controllers.MainLayoutTemplate().ParseFiles(
			"web/views/publish/layouts/central-card.html",
			"web/views/publish/identity/managed/request-password-reset-success.html"))
}

type requestPasswordResetPresenter struct {
	Email        string `form:"email"`
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
	return controllers.Render(ctx, requestPasswordResetTemplate, nil)
}

func (c RequestPasswordResetController) onPost(ctx *fiber.Ctx) error {
	// e := ctx.FormValue("email", "")

	p := &requestPasswordResetPresenter{}
	err := ctx.BodyParser(p)
	if err != nil {
		p.ErrorMessage = "Invalid email. Check and try again."
		return controllers.Render(ctx, requestPasswordResetTemplate, p)
	}

	email, err := user.NewEmail(p.Email)
	if err != nil {
		p.ErrorMessage = fmt.Sprintf("%s is not a valid email. Check and try agian.", p.Email)
		return controllers.Render(ctx, requestPasswordResetTemplate, p)
	}

	err = c.interactor.RequestResetPassword(email)
	if err != nil {

		if errors.Is(err, application.ErrEntityNotFound) {
			p.ErrorMessage = "There is no user who signed up with that email. Check and try agian."
			return controllers.Render(ctx, requestPasswordResetTemplate, p)
		}
		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
			p.ErrorMessage = "Your email has not been verified by our system. Click on a verification link sent to your email then try resetting your password again."
			return controllers.Render(ctx, requestPasswordResetTemplate, p)

		}
		p.ErrorMessage = "We're unable to reset your password at the moment. Try again later."
		return controllers.Render(ctx, requestPasswordResetTemplate, p)
	}
	return controllers.Render(ctx, requestPasswordResetSuccessTemplate, nil)
}
