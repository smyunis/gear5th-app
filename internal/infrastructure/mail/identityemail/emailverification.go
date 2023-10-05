package identityemail

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/mail"
)

var emailVerificationTemplate *template.Template

//go:embed emailverification.html
var emailVerificationTemplateFile string

func init() {
	emailVerificationTemplate = template.Must(
		mail.EmailMainLayoutTemplate().Parse(emailVerificationTemplateFile))
}

type verificationEmailPresenter struct {
	Link string
}

type VerificationEmailSender struct {
	appURL             *url.URL
	digitalSignService application.DigitalSignatureService
	mailSender         mail.MailSender
	logger             application.Logger
}

func NewVerifcationEmailSender(config infrastructure.ConfigurationProvider,
	digitalSignService application.DigitalSignatureService,
	mailSender mail.MailSender,
	logger application.Logger) VerificationEmailSender {
	appurlstr := config.Get("APP_URL", "https://gear5th.com")
	a, err := url.Parse(appurlstr)
	if err != nil {
		panic(err.Error())
	}

	return VerificationEmailSender{
		a,
		digitalSignService,
		mailSender,
		logger,
	}
}

func (s VerificationEmailSender) SendMail(u user.User) error {

	token, err := s.digitalSignService.Generate(u.ID.String())
	if err != nil {
		return err
	}

	verificationURL := s.buildEmailVerificationURL(u.ID.String(), token)

	p := verificationEmailPresenter{
		Link: verificationURL,
	}

	var htmlStringBuilder strings.Builder
	err = emailVerificationTemplate.ExecuteTemplate(&htmlStringBuilder, "email-main", p)
	if err != nil {
		s.logger.Error("mail/emailverification/parse-html", err)
		return err
	}

	msg := mail.Message{
		From: "no-reply@localhost",
		To:   u.Email.String(),
		Headers: map[string]string{
			// "From":         "gear5th Advertising",
			"To":           u.Email.String(),
			"Content-Type": "text/html",
			"Subject":      "Verify your email",
		},
		Body: htmlStringBuilder.String(),
	}

	err = s.mailSender.SendMail(msg)
	if err != nil {
		s.logger.Error("mail/emailverification/send", err)
		return err
	}

	s.logger.Info("mail/verificationemail", fmt.Sprintf("Sending Verification Email to %s <-> %s\n", u.Email.String(), verificationURL))
	return nil

}

func (s VerificationEmailSender) buildEmailVerificationURL(signedUpUserID, token string) string {
	// <APP_URL>/publish/identity/managed/{userId}/verify-email?token={token}
	s.appURL.Path = fmt.Sprintf("/publish/identity/managed/%s/verify-email", signedUpUserID)
	q := s.appURL.Query()
	q.Set("token", token)
	s.appURL.RawQuery = q.Encode()
	return s.appURL.String()
}
