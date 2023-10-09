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

//go:embed requestpasswordreset.html
var requestPasswordResetTemplateFile string

var requestPasswordResetTemplate *template.Template

func init() {
	requestPasswordResetTemplate = template.Must(
		mail.EmailMainLayoutTemplate().Parse(requestPasswordResetTemplateFile))
}

type RequestPasswordResetEmailService struct {
	webURL             *url.URL
	mailSender         mail.MailSender
	digitalSignService application.DigitalSignatureService
	logger             application.Logger
}

func NewRequestPasswordResetEmailService(
	config infrastructure.ConfigurationProvider,
	mailSender mail.MailSender,
	digitalSignService application.DigitalSignatureService,
	logger application.Logger) RequestPasswordResetEmailService {
	weburlstr := config.Get("APP_URL", "https://gear5th.com")
	w, err := url.Parse(weburlstr)
	if err != nil {
		panic(err.Error())
	}
	return RequestPasswordResetEmailService{
		w,
		mailSender,
		digitalSignService,
		logger,
	}
}

func (s RequestPasswordResetEmailService) SendMail(u user.User, passwordResetToken string) error {

	passwordResetURL := s.buildPasswordResetWebURL(u, passwordResetToken)

	var htmlStringBuilder strings.Builder
	err := requestPasswordResetTemplate.ExecuteTemplate(&htmlStringBuilder, "email-main", passwordResetURL)
	if err != nil {
		s.logger.Error("mail/reset-password/parse-html", err)
		return err
	}

	msg := mail.Message{
		From: "no-reply@localhost",
		To:   u.Email.String(),
		Headers: map[string]string{
			"From":         "gear5th Advertising <no-reply@gear5th.com>",
			"To":           u.Email.String(),
			"Content-Type": "text/html",
			"Subject":      "Reset your password",
		},
		Body: htmlStringBuilder.String(),
	}

	err = s.mailSender.SendMail(msg)
	if err != nil {
		s.logger.Error("mail/reset-password/send", err)
		return err
	}

	s.logger.Info("mail/resetpassword", fmt.Sprintf("Sent reset password mail to %s <-> %s \n", u.Email.String(), passwordResetURL))

	return nil
}

func (r RequestPasswordResetEmailService) buildPasswordResetWebURL(u user.User, passwordResetToken string) string {
	// <APP_URL>/publish/identity/managed/{userID}/reset-password?token={passwordResetToken}
	r.webURL.Path = fmt.Sprintf("/publish/identity/managed/%s/reset-password", u.ID.String())
	q := r.webURL.Query()
	q.Set("token", passwordResetToken)
	r.webURL.RawQuery = q.Encode()
	return r.webURL.String()
}
