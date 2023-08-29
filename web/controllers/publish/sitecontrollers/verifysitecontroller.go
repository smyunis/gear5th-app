package sitecontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/siteinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var verifySiteTemplate *template.Template
var verifySiteResultTemplate *template.Template

func init() {
	verifySiteTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/sites/verify-site.html"))
	verifySiteResultTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/sites/verify-site-result.html"))
}

type verifyitePresenter struct {
	Nav                    string
	ErrorMessage           string
	SiteAdsTxtRecord       string
	SiteID                 string
	SiteDomain             string
	VerificationSuccessful bool
}

type VerifySiteController struct {
	authMiddleware middlewares.JwtAuthenticationMiddleware
	interactor     siteinteractors.SiteInteractor
	logger         application.Logger
}

func NewVerifySiteController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	interactor siteinteractors.SiteInteractor,
	logger application.Logger) VerifySiteController {
	return VerifySiteController{
		authMiddleware,
		interactor,
		logger,
	}
}

func (c *VerifySiteController) AddRoutes(router *fiber.Router) {
	(*router).Use("/sites", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/sites/verify-site-prompt/:siteId", c.verifyPrompt)
	(*router).Add(fiber.MethodGet, "/sites/verify-site/:siteId", c.verifySite)
}

func (c *VerifySiteController) verifyPrompt(ctx *fiber.Ctx) error {
	siteID := ctx.Params("siteId", "")
	siteDomain := ctx.Query("site-domain", "")
	if siteID == "" || siteDomain == "" {
		return ctx.Redirect("/pages/error.html")
	}

	p := &verifyitePresenter{
		Nav:        "sites",
		SiteID:     siteID,
		SiteDomain: siteDomain,
	}
	return controllers.Render(ctx, verifySiteTemplate, p)
}

func (c *VerifySiteController) verifySite(ctx *fiber.Ctx) error {
	siteID := ctx.Params("siteId", "")
	if siteID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	p := &verifyitePresenter{
		Nav: "sites",
	}

	err := c.interactor.VerifySite(shared.ID(siteID))
	if err != nil {
		c.logger.Error("site/verifysite", err)
		p.VerificationSuccessful = false
		return controllers.Render(ctx, verifySiteResultTemplate, p)

	}
	p.VerificationSuccessful = true
	return controllers.Render(ctx, verifySiteResultTemplate, p)
}
