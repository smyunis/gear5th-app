package sitecontrollers

import (
	"html/template"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/siteinteractors"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var createSiteTemplate *template.Template

func init() {
	createSiteTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/sites/create-site.html"))
}

type createSitePresenter struct {
	Nav          string
	ErrorMessage string
	SiteURL      string `form:"url"`
}

type CreateSiteController struct {
	authMiddleware middlewares.JwtAuthenticationMiddleware
	interactor     siteinteractors.SiteInteractor
	logger         application.Logger
}

func NewCreateSiteController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	interactor siteinteractors.SiteInteractor,
	logger application.Logger) CreateSiteController {
	return CreateSiteController{
		authMiddleware,
		interactor,
		logger,
	}
}

func (c *CreateSiteController) AddRoutes(router *fiber.Router) {
	(*router).Use("/sites", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/sites/create-site", c.createSiteonGet)
	(*router).Add(fiber.MethodPost, "/sites/create-site", c.createSiteonPost)
}

func (c *CreateSiteController) createSiteonGet(ctx *fiber.Ctx) error {
	return controllers.Render(ctx, createSiteTemplate, createSitePresenter{Nav: "sites"})
}

func (c *CreateSiteController) createSiteonPost(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	p := &createSitePresenter{
		Nav: "sites",
	}
	err = ctx.BodyParser(p)
	if err != nil {
		p.ErrorMessage = "One or more invalid inputs. Check and try again."
		return controllers.Render(ctx, createSiteTemplate, p)
	}

	u, err := url.Parse(p.SiteURL)
	if err != nil {
		p.ErrorMessage = "Invalid site URL. Check and try again."
		return controllers.Render(ctx, createSiteTemplate, p)
	}
	u.Fragment = ""
	u.Path = ""
	u.RawQuery = ""
	err = c.interactor.CreateSite(publisherID, *u)
	if err != nil {
		c.logger.Error("site/createsite", err)
		p.ErrorMessage = "We're unable to register your new site at the moment. Try again later."
		return controllers.Render(ctx, createSiteTemplate, p)
	}
	return ctx.Redirect("/publish/sites")
}
