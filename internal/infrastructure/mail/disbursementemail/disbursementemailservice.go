package disbursementemail

import (
	"fmt"
	"net/url"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/domain/payment/disbursement"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

type DisbursementEmailService struct {
	appURL             *url.URL
	digitalSignService application.DigitalSignatureService
	logger             application.Logger
}

func NewDisbursementEmailService(config infrastructure.ConfigurationProvider,
	digitalSignService application.DigitalSignatureService,
	logger application.Logger) DisbursementEmailService {
	appurlstr := config.Get("APP_URL", "https://gear5th.com")
	a, err := url.Parse(appurlstr)
	if err != nil {
		panic(err.Error())
	}

	return DisbursementEmailService{
		a,
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

	s.logger.Info("mail/disbursement", fmt.Sprintf("Send confirm disbursement Email to %s <-> %s\n", email.String(), confirmURL))
	s.logger.Info("mail/disbursement", fmt.Sprintf("Send reject disbursement Email to %s <-> %s\n", email.String(), rejectURL))
	return nil
}

func (s DisbursementEmailService) buildConfirmURL(disbursementID, token string) string {
	// <APP_URL>/payment/disbursement/{disbursementId}/confirm?token={token}
	s.appURL.Path = fmt.Sprintf("/payment/disbursement/%s/confirm", disbursementID)
	q := s.appURL.Query()
	q.Set("token", token)
	s.appURL.RawQuery = q.Encode()
	return s.appURL.String()
}

func (s DisbursementEmailService) buildRejectURL(disbursementID, token string) string {
	// <APP_URL>/payment/disbursement/{disbursementId}/reject?token={token}
	s.appURL.Path = fmt.Sprintf("/payment/disbursement/%s/reject", disbursementID)
	q := s.appURL.Query()
	q.Set("token", token)
	s.appURL.RawQuery = q.Encode()
	return s.appURL.String()
}
