package homecontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

var homeTemplate *template.Template

func init() {
	homeTemplate = template.Must(
		controllers.MainLayoutTemplate().ParseFiles(
			"web/views/publish/home/home.html"))

}

type HomeController struct{}

func NewHomeController() HomeController {
	return HomeController{}
}

func (c *HomeController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/home", c.onGet)
	(*router).Add(fiber.MethodPost, "/home", c.onPost)
}

func (c *HomeController) onGet(ctx *fiber.Ctx) error {
	return controllers.Render(ctx, homeTemplate, nil)
}

func (c *HomeController) onPost(ctx *fiber.Ctx) error {
	return nil
}
