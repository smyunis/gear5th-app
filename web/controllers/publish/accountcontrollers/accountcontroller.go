package accountcontrollers

import (
	"errors"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var accountTemplate *template.Template
var changePasswordTemplate *template.Template

func init() {
	accountTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/account/account.html"))
	changePasswordTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/account/change-password.html"))
}

type accountPresenter struct {
	IsManagedUser bool
	Email         string `form:"email"`
	FirstName     string `form:"first-name"`
	LastName      string `form:"last-name"`
	PhoneNumber   string `form:"phone-number"`
	ErrorMessage  string
}

type AccountController struct {
	authMiddleware middlewares.JwtAuthenticationMiddleware
	interactor     identityinteractors.UserAccountInteractor
	logger         application.Logger
}

func NewAccountController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	interactor identityinteractors.UserAccountInteractor,
	logger application.Logger) AccountController {
	return AccountController{
		authMiddleware,
		interactor,
		logger,
	}
}

func (c *AccountController) AddRoutes(router *fiber.Router) {
	(*router).Use("/account", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/account", c.accountOnGet)
	(*router).Add(fiber.MethodPost, "/account", c.accountOnPost)
	(*router).Add(fiber.MethodGet, "/account/change-password", c.changePasswordOnGet)
	(*router).Add(fiber.MethodPost, "/account/change-password", c.changePasswordOnPost)
}

func (c *AccountController) accountOnGet(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	u, err := c.interactor.User(publisherID)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	p := accountPresenter{
		IsManagedUser: u.AuthenticationMethod == user.Managed,
		Email:         u.Email.String(),
		PhoneNumber:   u.PhoneNumber.String(),
		FirstName:     u.Fullname.FirstName(),
		LastName:      u.Fullname.LastName(),
	}

	return controllers.Render(ctx, accountTemplate, p)
}

func (c *AccountController) accountOnPost(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	p := &accountPresenter{}
	err = ctx.BodyParser(p)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	phone, err := user.NewPhoneNumber(p.PhoneNumber)
	if err != nil {
		p.ErrorMessage = "Invalid phone number. Check and try again."
		return controllers.Render(ctx, accountTemplate, p)
	}

	profile := identityinteractors.UserProfile{
		PhoneNumber: phone,
		Fullname:    user.NewPersonName(p.FirstName, p.LastName),
	}

	err = c.interactor.SetUser(publisherID, profile)
	if err != nil {
		p.ErrorMessage = "We're unable to update your account profile at the moment. Try again later."
		return controllers.Render(ctx, accountTemplate, p)
	}

	return ctx.Redirect("/publish/account")
}

func (c *AccountController) changePasswordOnGet(ctx *fiber.Ctx) error {
	return controllers.Render(ctx, changePasswordTemplate, nil)
}

type changePasswordPresenter struct {
	CurrentPassword string `form:"current-password"`
	NewPassword     string `form:"password"`
	ErrorMessage    string
}

func (c *AccountController) changePasswordOnPost(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	p := &changePasswordPresenter{}
	err = ctx.BodyParser(p)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}

	err = c.interactor.ChangePassword(publisherID, p.CurrentPassword, p.NewPassword)
	if err != nil {
		if errors.Is(err, application.ErrAuthorization) {
			p.ErrorMessage = "Seems like your current password is incorrect. Check and try again."
			return controllers.Render(ctx, changePasswordTemplate, p)
		}
		p.ErrorMessage = "We're unable to change your password at the moment. Try again later."
		return controllers.Render(ctx, changePasswordTemplate, p)
	}

	return ctx.Redirect("/publish/identity/signin")
}
