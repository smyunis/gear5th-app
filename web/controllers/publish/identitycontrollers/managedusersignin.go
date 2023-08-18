package identitycontrollers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
)

type ManagedUserController struct {
	// controllers.Controller
	interactor manageduserinteractors.ManagedUserInteractor
}

func NewManagedUserController(interactor manageduserinteractors.ManagedUserInteractor) ManagedUserController {
	return ManagedUserController{
		interactor,
	}
}

// GET
func (c ManagedUserController) SignIn(ctx *fiber.Ctx) error {

	return ctx.Render("publish/identity/signin", nil, "publish/layouts/main")
}

type MangedUserSigninPresenter struct {
	HasValidationErrors bool
}

// POST
func (c ManagedUserController) ManagedUserSignIn(ctx *fiber.Ctx) error {

	userEmail := ctx.FormValue("email", "")
	password := ctx.FormValue("password", "")

	if userEmail == "" || password == "" {
		// ctx.SendStatus(fiber.StatusUnauthorized)
		return ctx.Render("publish/identity/signin", MangedUserSigninPresenter{
			HasValidationErrors: true,
		}, "publish/layouts/main")

	}

	return ctx.Render("publish/identity/signin", MangedUserSigninPresenter{
		HasValidationErrors: true,
	}, "publish/layouts/main")

	// 	email, err := user.NewEmail(userEmail)
	// 	if err != nil {
	// 		return c.SendProblemDetails(ctx, fiber.StatusBadRequest,
	// 			"Invalid Email",
	// 			fmt.Sprintf("%s is not a valid email", email.String()))
	// 	}

	// 	token, err := c.interactor.SignIn(email, password)
	// 	if err != nil {
	// 		if errors.Is(err, identityinteractors.ErrEmailNotVerified) {
	// 			return c.SendProblemDetails(ctx, fiber.StatusPreconditionRequired,
	// 				"Unverified Email",
	// 				fmt.Sprintf("email %s is not verified", email.String()))
	// 		}
	// 		return c.SendProblemDetails(ctx, fiber.StatusUnauthorized, "", "")
	// 	}

	//		return ctx.JSON(struct {
	//			AccessToken string `json:"accessToken"`
	//		}{token})
	//	}

}
