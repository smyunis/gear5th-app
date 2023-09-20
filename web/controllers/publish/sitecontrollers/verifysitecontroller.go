package sitecontrollers

import (
	"errors"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
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
	ErrorMessage           string
	SiteAdsTxtRecord       string
	SiteID                 string
	SiteDomain             string
	VerificationSuccessful bool
}

type VerifySiteController struct {
	authMiddleware middlewares.JwtAuthenticationMiddleware
	interactor     publisherinteractors.SiteInteractor
	logger         application.Logger
}

func NewVerifySiteController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	interactor publisherinteractors.SiteInteractor,
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
	if siteID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	s, err := c.interactor.Site(shared.ID(siteID))
	if err != nil {
		if !errors.Is(err, application.ErrEntityNotFound) {
			c.logger.Error("site/getverifyprompt", err)
		}
		return ctx.Redirect("/pages/error.html")
	}

	siteURL := s.URL
	adsTxtRecord, err := c.interactor.GenerateAdsTxtRecord(shared.ID(siteID))
	if err != nil {
		if !errors.Is(err, application.ErrEntityNotFound) {
			c.logger.Error("site/getverifyprompt", err)
		}
		return ctx.Redirect("/pages/error.html")
	}
	p := &verifyitePresenter{
		SiteID:           siteID,
		SiteDomain:       siteURL.String(),
		SiteAdsTxtRecord: adsTxtRecord.String(),
	}
	return controllers.Render(ctx, verifySiteTemplate, p)
}

func (c *VerifySiteController) verifySite(ctx *fiber.Ctx) error {
	siteID := ctx.Params("siteId", "")
	if siteID == "" {
		return ctx.Redirect("/pages/error.html")
	}

	p := &verifyitePresenter{
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
