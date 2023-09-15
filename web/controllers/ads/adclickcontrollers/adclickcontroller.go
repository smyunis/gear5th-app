package adclickcontrollers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/adsinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/application/advertiserinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type AdClickController struct {
	adPieceInteractor  advertiserinteractors.AdPieceInteractor
	adClickInteractor  adsinteractors.AdClickInteractor
	logger  application.Logger
}

func NewAdClickController(
	adPieceInteractor  advertiserinteractors.AdPieceInteractor,
	adClickInteractor  adsinteractors.AdClickInteractor,
	logger application.Logger) AdClickController {
	return AdClickController{
		adPieceInteractor,
		adClickInteractor,
		logger,
	}
}

func (c *AdClickController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/adclick/adpiece/:adPieceId", c.referralOnGet)
}

func (c *AdClickController) referralOnGet(ctx *fiber.Ctx) error {
	//TODO count ad click here
	adPieceID := ctx.Params("adPieceId", "")
	if adPieceID == "" {
		return nil
	}

	token := ctx.Query("token", "")
	if token == "" {
		return nil
	}

	err := c.adClickInteractor.OnClick(shared.ID(adPieceID), token)
	if err != nil {
		c.logger.Error("adpiece/adclick", err)
	}

	a, err := c.adPieceInteractor.AdPiece(shared.ID(adPieceID))
	if err != nil {
		c.logger.Error("adpiece/get", err)
		return nil
	}

	if a.Ref.String() == "#" || a.Ref.String() == "" {
		return nil
	}

	return ctx.Redirect(a.Ref.String())
}