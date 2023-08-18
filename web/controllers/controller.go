package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/pkg/problemdetails"
)

type ControllerRoute struct {
	Method string
	Path   string
	Handler func(*fiber.Ctx) error
}
type ControllerRoutes []ControllerRoute

type Controller struct {
	// Handlers ControllerHandlers
	Method string
	Path   string
}

func (Controller) SendProblemDetails(ctx *fiber.Ctx, status int, title, detail string) error {
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
