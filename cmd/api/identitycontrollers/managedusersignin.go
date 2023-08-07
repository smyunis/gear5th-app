package identitycontrollers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application/identity/usersignin"
)

type ManagedUserSignInController struct {
	interactor usersignin.ManagedUserInteractor
}

func NewManagedUserSignIn(interactor usersignin.ManagedUserInteractor) ManagedUserSignInController {
	return ManagedUserSignInController{
		interactor,
	}
}

func (c ManagedUserSignInController) SignIn(ctx *fiber.Ctx) error {

	return ctx.SendString("Alhamduli Allah")
}
