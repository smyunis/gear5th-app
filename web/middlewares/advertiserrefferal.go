package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

type AdvertiserRefferalMiddleware struct {
	// accessTokenService identityinteractors.AccessTokenService
	digitalSignatureService application.DigitalSignatureService
}

func NewAdvertiserRefferalMiddleware(
	// accessTokenService identityinteractors.AccessTokenService,
	digitalSignatureService application.DigitalSignatureService) AdvertiserRefferalMiddleware {
	return AdvertiserRefferalMiddleware{
		// accessTokenService,
		digitalSignatureService,
	}
}

func (m *AdvertiserRefferalMiddleware) Verification(ctx *fiber.Ctx) error {

	token := ctx.Query("token", "")
	if token == "" {
		return ctx.Redirect("/pages/error.html")
	}

	if !m.digitalSignatureService.Validate(token) {
		return ctx.Redirect("/pages/error.html")
	}

	actorUserID, err := m.digitalSignatureService.GetMessage(token)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	ctx.Locals(controllers.ActorUserID, shared.ID(actorUserID))
	ctx.Locals(controllers.AdvertiserToken, token)
	return ctx.Next()
}

func (m *AdvertiserRefferalMiddleware) ActorUserID(ctx *fiber.Ctx) (shared.ID, error) {
	if id, ok := ctx.Locals(controllers.ActorUserID).(shared.ID); ok {
		return id, nil
	}
	return shared.ID("err"), ErrInvalidActorID
}

func (m *AdvertiserRefferalMiddleware) AdvertiserToken(ctx *fiber.Ctx) (string, error) {
	if id, ok := ctx.Locals(controllers.AdvertiserToken).(string); ok {
		return id, nil
	}
	return "", ErrInvalidActorID
}
