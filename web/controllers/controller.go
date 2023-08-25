package controllers

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/pkg/problemdetails"
)

func CardMainLayoutTemplate() *template.Template {
	mainTmpl := template.Must(template.ParseFiles("web/views/publish/layouts/card-main.html"))
	return template.Must(mainTmpl.ParseGlob("web/views/publish/layouts/components/*.html"))
}

func ConsoleMainLayoutTemplate() *template.Template {
	mainTmpl := template.Must(template.ParseFiles("web/views/publish/layouts/console-main.html"))
	return template.Must(mainTmpl.ParseGlob("web/views/publish/layouts/components/*.html"))
}

func Render(ctx *fiber.Ctx, t *template.Template, binding any) error {
	ctx.Set("Content-Type", "text/html")
	return t.Execute(ctx.Response().BodyWriter(), binding)
}

func SendProblemDetails(ctx *fiber.Ctx, status int, title, detail string) error {
	prob := problemdetails.NewProblemDetails(status)
	if title != "" {
		prob.Title = title
	}
	if detail != "" {
		prob.Detail = detail
	}
	ctx.SendStatus(status)
	return ctx.JSON(prob)
}
