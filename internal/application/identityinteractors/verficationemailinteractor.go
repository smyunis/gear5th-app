package identityinteractors

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
)

type VerificationEmailService interface {
	SendMail(u user.User) error
}

type VerificationEmailInteractor struct {
	emailService VerificationEmailService
	logger       application.Logger
}

func NewVerificationEmailInteractor(
	emailService VerificationEmailService,
	logger application.Logger) VerificationEmailInteractor {
	return VerificationEmailInteractor{
		emailService,
		logger,
	}
}

func (i *VerificationEmailInteractor) HandleUserSignedUpEvent(event any) {
	u := event.(user.User)
	err := i.emailService.SendMail(u)
	if err != nil {
		i.logger.Error("publisher/signup/verificationemail", err)
	}
}
