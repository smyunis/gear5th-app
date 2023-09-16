package identityemail

import (
	"fmt"
	"net/url"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

type VerifcationEmailSender struct {
	appURL             *url.URL
	digitalSignService application.DigitalSignatureService
	logger             application.Logger
}

func NewVerifcationEmailSender(config infrastructure.ConfigurationProvider,
	digitalSignService application.DigitalSignatureService,
	logger application.Logger) VerifcationEmailSender {
	appurlstr := config.Get("APP_URL", "https://gear5th.com")
	a, err := url.Parse(appurlstr)
	if err != nil {
		panic(err.Error())
	}

	return VerifcationEmailSender{
		a,
		digitalSignService,
		logger,
	}
}

func (s VerifcationEmailSender) SendMail(u user.User) error {

	token, err := s.digitalSignService.Generate(u.ID.String())
	if err != nil {
		return err
	}

	verificationURL := s.buildEmailVerificationURL(u.ID.String(), token)

	//TODO send email with link to verify email

	s.logger.Info("mail/verificationemail", fmt.Sprintf("Sending Verification Email to %s <-> %s\n", u.Email.String(), verificationURL))
	return nil

}

func (s VerifcationEmailSender) buildEmailVerificationURL(signedUpUserID, token string) string {
	// <APP_URL>/publish/identity/managed/{userId}/verify-email?token={token}
	s.appURL.Path = fmt.Sprintf("/publish/identity/managed/%s/verify-email", signedUpUserID)
	q := s.appURL.Query()
	q.Set("token", token)
	s.appURL.RawQuery = q.Encode()
	return s.appURL.String()
}
