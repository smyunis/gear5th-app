package homecontrollers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
	"gitlab.com/gear5th/gear5th-app/web/middlewares"
)

var homeTemplate *template.Template

func init() {
	homeTemplate = template.Must(
		controllers.ConsoleMainLayoutTemplate().ParseFiles(
			"web/views/publish/home/home.html"))
}

type homePresenter struct {
}

type HomeController struct {
	authMiddleware middlewares.JwtAuthenticationMiddleware
}

func NewHomeController(authMiddleware middlewares.JwtAuthenticationMiddleware) HomeController {
	return HomeController{
		authMiddleware,
	}
}

func (c *HomeController) AddRoutes(router *fiber.Router) {
	(*router).Use("/home", c.authMiddleware.Authentication)
	(*router).Add(fiber.MethodGet, "/home", c.onGet)
}

func (c *HomeController) onGet(ctx *fiber.Ctx) error {

	// actorUserID := ctx.Locals(controllers.ActorUserID).(shared.ID)

	p := homePresenter{
	}

	return controllers.Render(ctx, homeTemplate, p)
}

func (c *HomeController) onPost(ctx *fiber.Ctx) error {
	return nil
}
