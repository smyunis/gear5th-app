package publishercontrollers

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/web/controllers"
)

var publisherSignUpTemplate *template.Template
var publisherSignUpSuccessTemplate *template.Template

func init() {
	publisherSignUpTemplate = template.Must(controllers.CardMainLayoutTemplate().ParseFiles(
		"web/views/publish/identity/signup.html"))

	publisherSignUpSuccessTemplate = template.Must(controllers.CardMainLayoutTemplate().ParseFiles(
		"web/views/publish/identity/signup-success.html"))
}

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
	interactor publisherinteractors.PublisherInteractor
	logger     application.Logger
}

func NewPublisherSignUpController(interactor publisherinteractors.PublisherInteractor,
	logger application.Logger) PublisherSignUpController {
	return PublisherSignUpController{
		interactor,
		logger,
	}
}

func (c *PublisherSignUpController) AddRoutes(router *fiber.Router) {
	(*router).Add(fiber.MethodGet, "/identity/signup", c.publisherSignUpOnGet)
	(*router).Add(fiber.MethodPost, "/identity/signup", c.publisherSignUpOnPost)
	(*router).Add(fiber.MethodPost, "/identity/oauth/signup", c.oauthUserPublisherSignUp)
}

func (c *PublisherSignUpController) publisherSignUpOnGet(ctx *fiber.Ctx) error {
	//TODO set cacheing headers
	return controllers.Render(ctx, publisherSignUpTemplate, nil)
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
		return controllers.Render(ctx, publisherSignUpTemplate, p)
	}

	email, err := user.NewEmail(pub.Email)
	if err != nil {
		p.InvalidEmail = true
		p.ErrorMessage = fmt.Sprintf("%s is not a valid email. Check and try again.", pub.Email)
		return controllers.Render(ctx, publisherSignUpTemplate, p)
	}

	u := user.NewUser(email)

	if pub.PhoneNumber != "" {
		uPhoneNum, err := user.NewPhoneNumber(pub.PhoneNumber)
		if err != nil {
			p.InvalidPhoneNumber = true
			p.ErrorMessage = fmt.Sprintf("%s is not a valid phone number. Check and try again.", pub.PhoneNumber)
			return controllers.Render(ctx, publisherSignUpTemplate, p)
		}
		u.PhoneNumber = uPhoneNum
	}

	name := user.NewPersonName(pub.FirstName, pub.LastName)
	mu := u.AsManagedUser(name, pub.Password)

	err = c.interactor.ManagedUserSignUp(u, mu)
	if err != nil {
		if errors.Is(err, application.ErrConflictFound) {
			p.InvalidEmail = true
			p.ErrorMessage = "It seems another person has sign up with this email already. Check and try again."
			return controllers.Render(ctx, publisherSignUpTemplate, p)
		}

		c.logger.Error("publishers/signup", err)
		p.ErrorMessage = "We're unable to sign you up at the moment. Try again later."
		return controllers.Render(ctx, publisherSignUpTemplate, p)
	}

	return controllers.Render(ctx, publisherSignUpSuccessTemplate, nil)
}

type oAuthSignUpPresenter struct {
	ClientID   string `form:"clientId"`
	Credential string `form:"credential"`
	SelectBy   string `form:"select_by"`
	CSRFToken  string `form:"g_csrf_token"`
}

func (c *PublisherSignUpController) oauthUserPublisherSignUp(ctx *fiber.Ctx) error {
	p := &oAuthSignUpPresenter{}
	err := ctx.BodyParser(p)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	err = c.interactor.OAuthUserSignUp(p.Credential)
	if err != nil {
		return ctx.Redirect("/pages/error.html")
	}
	return controllers.Render(ctx, publisherSignUpSuccessTemplate, nil)
}
