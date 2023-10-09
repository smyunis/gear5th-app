package disbursementemail

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/disbursement"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/mail"
)

//go:embed requestdisbursement.html
var requestDisbursementTemplateFile string

var requestDisbursementTemplate *template.Template

func init() {
	requestDisbursementTemplate = template.Must(
		mail.EmailMainLayoutTemplate().Parse(requestDisbursementTemplateFile))
}

type disbursementRequestPresenter struct {
	ConfirmLink  string
	RejectLink   string
	Disbursement disbursement.Disbursement
}

type DisbursementEmailService struct {
	appURL             *url.URL
	mailSender         mail.MailSender
	digitalSignService application.DigitalSignatureService
	logger             application.Logger
}

func NewDisbursementEmailService(config infrastructure.ConfigurationProvider,
	mailSender mail.MailSender,
	digitalSignService application.DigitalSignatureService,
	logger application.Logger) DisbursementEmailService {
	appurlstr := config.Get("APP_URL", "https://gear5th.com")
	a, err := url.Parse(appurlstr)
	if err != nil {
		panic(err.Error())
	}

	return DisbursementEmailService{
		a,
		mailSender,
		digitalSignService,
		logger,
	}
}

func (s DisbursementEmailService) SendRequestDisbursementConfirmation(email user.Email, d disbursement.Disbursement) error {

	token, err := s.digitalSignService.Generate(d.ID.String())
	if err != nil {
		return err
	}

	confirmURL := s.buildConfirmURL(d.ID.String(), token)
	rejectURL := s.buildRejectURL(d.ID.String(), token)

	p := disbursementRequestPresenter{
		ConfirmLink: confirmURL,
		RejectLink:  rejectURL,
		Disbursement: d,
	}

	var htmlStringBuilder strings.Builder
	err = requestDisbursementTemplate.ExecuteTemplate(&htmlStringBuilder, "email-main", p)
	if err != nil {
		s.logger.Error("mail/request-disbursement/parse-html", err)
		return err
	}

	msg := mail.Message{
		From: "no-reply@localhost",
		To:   email.String(),
		Headers: map[string]string{
			"From":         "gear5th Advertising <no-reply@gear5th.com>",
			"To":           email.String(),
			"Content-Type": "text/html",
			"Subject":      "Confirm disbursement",
		},
		Body: htmlStringBuilder.String(),
	}

	err = s.mailSender.SendMail(msg)
	if err != nil {
		s.logger.Error("mail/request-disbursement/send", err)
		return err
	}

	s.logger.Info("mail/disbursement", fmt.Sprintf("Send confirm disbursement Email to %s <-> %s\n", email.String(), confirmURL))
	s.logger.Info("mail/disbursement", fmt.Sprintf("Send reject disbursement Email to %s <-> %s\n", email.String(), rejectURL))
	return nil
}

func (s DisbursementEmailService) buildConfirmURL(disbursementID, token string) string {
	// <APP_URL>/publish/payments/disbursement/{disbursementId}/confirm?token={token}
	s.appURL.Path = fmt.Sprintf("/publish/payments/disbursement/%s/confirm", disbursementID)
	q := s.appURL.Query()
	q.Set("token", token)
	s.appURL.RawQuery = q.Encode()
	return s.appURL.String()
}

func (s DisbursementEmailService) buildRejectURL(disbursementID, token string) string {
	// <APP_URL>/publish/payments/disbursement/{disbursementId}/reject?token={token}
	s.appURL.Path = fmt.Sprintf("/publish/payments/disbursement/%s/reject", disbursementID)
	q := s.appURL.Query()
	q.Set("token", token)
	s.appURL.RawQuery = q.Encode()
	return s.appURL.String()
}
