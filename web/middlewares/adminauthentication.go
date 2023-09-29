package middlewares

import (

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)


type AdminAuthenticationMiddleware struct {
	accessTokenService application.AccessTokenService
}

func NewAdminAuthenticationMiddleware(
	accessTokenService application.AccessTokenService) AdminAuthenticationMiddleware {
	return AdminAuthenticationMiddleware{
		accessTokenService,
	}
}

func (m *AdminAuthenticationMiddleware) Authentication(ctx *fiber.Ctx) error {
	accessToken := ctx.Cookies(controllers.AdminAccessTokenCookieName, "")
	if accessToken == "" {
		return ctx.Redirect("/admin/identity/signin")
	}
	if !m.accessTokenService.Validate(accessToken) {
		return ctx.Redirect("/admin/identity/signin")
	}
	actorUserID, err := m.accessTokenService.UserID(accessToken)
	if err != nil {
		return ctx.Redirect("/admin/identity/signin")
	}
	ctx.Locals(controllers.ActorUserID, shared.ID(actorUserID))
	return ctx.Next()
}

func (m *AdminAuthenticationMiddleware) ActorUserID(ctx *fiber.Ctx) (shared.ID, error) {
	if id, ok := ctx.Locals(controllers.ActorUserID).(shared.ID); ok {
		return id, nil
	}
	return shared.ID("err"), ErrInvalidActorID
}
