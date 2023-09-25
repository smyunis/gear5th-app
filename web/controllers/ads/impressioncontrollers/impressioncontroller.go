package impressioncontrollers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/advertiserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type impressionPersenter struct {
	AdPieceID   string `json:"adPieceId"`
	Token       string `json:"token"`
	SiteID      string `json:"siteId"`
	AdSlotID    string `json:"adSlotId"`
	PublihserID string `json:"publihserId"`
	Origin      string `json:"origin"`
}

type ImpressionController struct {
	adPieceInteractor    advertiserinteractors.AdPieceInteractor
	impressionInteractor adsinteractors.AdsInteractor
	logger               application.Logger
}

func NewImpressionController(
	adPieceInteractor advertiserinteractors.AdPieceInteractor,
	impressionInteractor adsinteractors.AdsInteractor,
	logger application.Logger) ImpressionController {
	return ImpressionController{
		adPieceInteractor,
		impressionInteractor,
		logger,
	}
}

func (c *ImpressionController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodPost, "/impression", c.impressionOnPost)
}

func (c *ImpressionController) impressionOnPost(ctx *fiber.Ctx) error {
	p := &impressionPersenter{}
	err := ctx.BodyParser(p)
	if err != nil {
		return nil
	}

	err = c.impressionInteractor.NewImpression(shared.ID(p.AdPieceID), shared.ID(p.SiteID), shared.ID(p.AdSlotID), shared.ID(p.PublihserID), p.Token, p.Origin)
	if err != nil {
		c.logger.Error("impression/recieved", err)
		return nil
	}

	return nil
}
