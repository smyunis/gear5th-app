package mail

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

//go:embed layout.html
var emailLayoutTemplateFile string

func EmailMainLayoutTemplate() *template.Template {
	return template.Must(template.New("email-main").Parse(emailLayoutTemplateFile))
}

type Message struct {
	From    string
	To      string
	Headers map[string]string
	Body    string
}

type MailSender interface {
	SendMail(email Message) error
}

type SMTPMailSender struct {
	host     string
	fromAddr string
	fromPass string
	logger   application.Logger
}

func NewSMTPMailSender(config infrastructure.ConfigurationProvider, logger application.Logger) SMTPMailSender {
	host := config.Get("SMTP_HOST", "")
	if host == "" {
		panic("config error, SMTP_HOST not found")
	}
	fa := config.Get("SMTP_SENDER", "")
	if fa == "" {
		panic("config error, SMTP_SENDER not found")
	}
	fp := config.Get("SMTP_SENDER_PASS", "")
	if fp == "" {
		panic("config error, SMTP_SENDER_PASS not found")
	}
	return SMTPMailSender{
		host,
		fa,
		fp,
		logger,
	}
}

func (s SMTPMailSender) SendMail(email Message) error {

	var emailHeadersSB strings.Builder
	for headerKey, headerValue := range email.Headers {
		emailHeadersSB.Write([]byte(fmt.Sprintf("%s: %s\r\n", headerKey, headerValue)))
	}

	msg := []byte(fmt.Sprintf("%s\r\n%s", emailHeadersSB.String(), email.Body))

	//TODO add sasl auth to postfix
	// auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(s.host, nil, email.From, []string{email.To}, msg)
	if err != nil {
		s.logger.Error("sending/mail", err)
		return err
	}
	return nil
}
