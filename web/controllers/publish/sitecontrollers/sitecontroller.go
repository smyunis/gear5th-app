package sitecontrollers

import (
	"errors"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/siteinteractors"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var siteTemplate *template.Template

func init() {
	siteTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/sites/sites.html"))
}

type sitePresenterActiveSite struct {
	SiteDomain     string
	IsSiteVerified bool
	SiteID         string
}
type sitePresenter struct {
	Nav          string
	ActiveSites  []sitePresenterActiveSite
	ErrorMessage string
}

type SiteController struct {
	authMiddleware middlewares.JwtAuthenticationMiddleware
	interactor     siteinteractors.SiteInteractor
	logger         application.Logger
}

func NewSiteController(authMiddleware middlewares.JwtAuthenticationMiddleware,
	interactor siteinteractors.SiteInteractor,
	logger application.Logger) SiteController {
	return SiteController{
		authMiddleware,
		interactor,
		logger,
	}
}

func (c *SiteController) AddRoutes(router *fiber.Router) {
	(*router).Use("/sites", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/sites", c.onGet)
}

func (c *SiteController) onGet(ctx *fiber.Ctx) error {
	publisherID, err := c.authMiddleware.ActorUserID(ctx)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	activeSites, err := c.interactor.ActiveSitesForPublisher(publisherID)
	if err != nil {
		if errors.Is(err, application.ErrEntityNotFound) {
			p := sitePresenter{
				Nav:         "sites",
				ActiveSites: make([]sitePresenterActiveSite, 0),
			}
			return controllers.Render(ctx, siteTemplate, p)
		}
		p := sitePresenter{
			Nav:          "sites",
			ErrorMessage: "We're unable to get your sites at the moment. Try again later.",
		}
		c.logger.Error("site/activesitesforpublisher", err)
		return controllers.Render(ctx, siteTemplate, p)
	}
	sites := make([]sitePresenterActiveSite, 0)
	for _, s := range activeSites {
		siteURL := s.URL()
		siteDomain := siteURL.Hostname()
		sites = append(sites, sitePresenterActiveSite{
			SiteID:         string(s.SiteID()),
			SiteDomain:     siteDomain,
			IsSiteVerified: s.IsVerified(),
		})
	}
	p := sitePresenter{
		Nav:         "sites",
		ActiveSites: sites,
	}
	return controllers.Render(ctx, siteTemplate, p)
}
