package publishercontrollers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/gear5th/gear5th-api/internal/application/publisherinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/pkg/problemdetails"
)

type PublisherSignUpController struct {
	Method     string
	Path       string
	interactor publisherinteractors.PublisherSignUpInteractor
}

func NewPublisherSignUpController(interactor publisherinteractors.PublisherSignUpInteractor) PublisherSignUpController {
	return PublisherSignUpController{
		fiber.MethodPost,
		"/publishers/managed",
		interactor,
	}
}

func (c *PublisherSignUpController) ManagedUserSignUp(ctx *fiber.Ctx) error {
	pub := &publisherManaged{}
	err := ctx.BodyParser(pub)
	if err != nil {
		return ctx.JSON(problemdetails.NewProblemDetails(fiber.StatusBadRequest))
	}

	email, err := user.NewEmail(pub.Email)
	if err != nil {
		probDetails := problemdetails.NewProblemDetails(fiber.StatusBadRequest)
		probDetails.Detail = fmt.Sprintf("%s is not a valid email", email.Email())
		return ctx.JSON(probDetails)
	}

	u := user.NewUser(email)

	if pub.PhoneNumber != "" {
		uPhoneNum, err := user.NewPhoneNumber(pub.PhoneNumber)
		if err != nil {
			probDetails := problemdetails.NewProblemDetails(fiber.StatusBadRequest)
			probDetails.Detail = fmt.Sprintf("%s is not a valid phone number", uPhoneNum.PhoneNumber())
			return ctx.JSON(probDetails)
		}
		u.SetPhoneNumber(uPhoneNum)
	}

	name := user.NewPersonName(pub.FirstName, pub.LastName)
	mu := u.AsManagedUser(name, pub.Password)

	err = c.interactor.ManagedUserSignUp(&u, &mu)
	if err != nil {
		return ctx.JSON(problemdetails.NewProblemDetails(fiber.StatusInternalServerError))
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

type publisherManaged struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}
