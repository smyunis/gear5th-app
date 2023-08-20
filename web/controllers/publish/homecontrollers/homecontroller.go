package homecontrollers

import "github.com/gofiber/fiber/v2"

type HomeController struct {
}

func NewHomeController() HomeController {
	return HomeController{}
}

func (c *HomeController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/home", c.homeOnGet)
	(*router).Add(fiber.MethodPost, "/home", c.homeOnPost)
}

func (c *HomeController) homeOnGet(ctx *fiber.Ctx) error {
	return nil
}

func (c *HomeController) homeOnPost(ctx *fiber.Ctx) error {
	return nil
}
