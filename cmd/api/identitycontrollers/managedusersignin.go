package identitycontrollers

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/cmd/api/problemdetails"
	"gitlab.com/gear5th/gear5th-api/internal/application/identity/usersignin"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type ManagedUserController struct {
	Method     string
	Path       string
	interactor usersignin.ManagedUserInteractor
}

func NewManagedUserController(interactor usersignin.ManagedUserInteractor) ManagedUserController {
	return ManagedUserController{
		fiber.MethodPost,
		"/identity/managed/signin",
		interactor,
	}
}

func (c ManagedUserController) SignIn(ctx *fiber.Ctx) error {

	// email is passed in as username
	username, password, err := c.basicAuthCredentials(ctx)
	if err != nil {
		return ctx.JSON(problemdetails.NewProblemDetails(fiber.StatusUnauthorized))
	}

	email, err := user.NewEmail(username)
	if err != nil {
		return ctx.JSON(problemdetails.NewProblemDetails(fiber.StatusUnauthorized))
	}

	token, err := c.interactor.SignIn(email, password)
	if err != nil {
		return ctx.JSON(problemdetails.NewProblemDetails(fiber.StatusUnauthorized))
	}

	return ctx.JSON(struct {
		AccessToken string `json:"accessToken"`
	}{token})
}

func (ManagedUserController) basicAuthCredentials(c *fiber.Ctx) (username, password string, err error) {

	auth := c.Get(fiber.HeaderAuthorization)

	// Check if the header contains content besides "basic".
	if len(auth) <= 6 || !strings.EqualFold(auth[:6], "basic ") {
		return "", "", fiber.ErrUnauthorized
	}

	// Decode the header contents
	raw, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return "", "", fiber.ErrUnauthorized
	}

	// Get the credentials
	creds := string(raw)

	// Check if the credentials are in the correct form
	// which is "username:password".
	index := strings.Index(creds, ":")
	if index == -1 {
		return "", "", fiber.ErrUnauthorized
	}

	return creds[:index], creds[index+1:], nil
}
