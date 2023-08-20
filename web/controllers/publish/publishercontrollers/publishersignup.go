package publishercontrollers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type publisherSignUp struct {
	FirstName          string `form:"first-name"`
	LastName           string `form:"last-name"`
	Email              string `form:"email"`
	Password           string `form:"password"`
	PhoneNumber        string `form:"phone-number"`
	ErrorMessage       string
	InvalidEmail       bool
	InvalidPhoneNumber bool
}

type PublisherSignUpController struct {
	interactor publisherinteractors.PublisherSignUpInteractor
}

func NewPublisherSignUpController(interactor publisherinteractors.PublisherSignUpInteractor) PublisherSignUpController {
	return PublisherSignUpController{
		interactor,
	}
}

func (c *PublisherSignUpController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/signup", c.publisherSignUpOnGet)
	(*router).Add(fiber.MethodPost, "/identity/signup", c.publisherSignUpOnPost)
}

func (c *PublisherSignUpController) publisherSignUpOnGet(ctx *fiber.Ctx) error {
	//TODO set cacheing headers
	return ctx.Render("publish/identity/signup", nil, "publish/layouts/main")
}

func (c *PublisherSignUpController) publisherSignUpOnPost(ctx *fiber.Ctx) error {
	pub := &publisherSignUp{}
	err := ctx.BodyParser(pub)
	p := publisherSignUp{
		Email:       pub.Email,
		FirstName:   pub.FirstName,
		LastName:    pub.LastName,
		PhoneNumber: pub.PhoneNumber,
		Password:    pub.Password,
	}
	if err != nil {
		p.ErrorMessage = "Unable to sign you up with the details you provided. Please review and try again."
		return ctx.Render("publish/identity/signup", p, "publish/layouts/main")
	}

	email, err := user.NewEmail(pub.Email)
	if err != nil {
		p.InvalidEmail = true
		p.ErrorMessage = fmt.Sprintf("%s is not a valid email. Check and try agian.", pub.Email)
		return ctx.Render("publish/identity/signup", p, "publish/layouts/main")
	}

	u := user.NewUser(email)

	if pub.PhoneNumber != "" {
		uPhoneNum, err := user.NewPhoneNumber(pub.PhoneNumber)
		if err != nil {
			p.InvalidPhoneNumber = true
			p.ErrorMessage = fmt.Sprintf("%s is not a valid phone number. Check and try agian.", pub.PhoneNumber)
			return ctx.Render("publish/identity/signup", p, "publish/layouts/main")
		}
		u.SetPhoneNumber(uPhoneNum)
	}

	name := user.NewPersonName(pub.FirstName, pub.LastName)
	mu := u.AsManagedUser(name, pub.Password)

	err = c.interactor.ManagedUserSignUp(u, mu)
	if err != nil {
		if errors.Is(err, application.ErrConflictFound) {
			p.InvalidEmail = true
			p.ErrorMessage = "It seems another person has sign up with this email already. Check and try again."
			return ctx.Render("publish/identity/signup", p, "publish/layouts/main")
		}

		p.ErrorMessage = "We're unable to sign you up at the moment. Try agian later."
		return ctx.Render("publish/identity/signup", p, "publish/layouts/main")
	}

	return ctx.Render("publish/identity/signup-success", nil, "publish/layouts/main")
}
