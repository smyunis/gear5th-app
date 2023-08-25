package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

type JwtAuthenticationMiddleware struct {
	accessTokenService identityinteractors.AccessTokenService
}

func NewJwtAuthenticationMiddleware(
	accessTokenService identityinteractors.AccessTokenService) JwtAuthenticationMiddleware {
	return JwtAuthenticationMiddleware{
		accessTokenService,
	}
}

func (m *JwtAuthenticationMiddleware) Authentication(ctx *fiber.Ctx) error {
	accessToken := ctx.Cookies(controllers.AccessTokenCookieName, "")
	if accessToken == "" {
		return ctx.Redirect("/publish/identity/signin", fiber.StatusUnauthorized)
	}
	if !m.accessTokenService.Validate(accessToken) {
		return ctx.Redirect("/publish/identity/signin", fiber.StatusUnauthorized)
	}
	actorUserID, err := m.accessTokenService.UserID(accessToken)
	if err != nil {
		return ctx.Redirect("/publish/identity/signin", fiber.StatusUnauthorized)
	}
	ctx.Locals(controllers.ActorUserID, shared.ID(actorUserID))
	return ctx.Next()
}
