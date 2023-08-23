package identityemail

import (
	"fmt"
	"net/url"

	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/application/identityinteractors"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

type RequestPassordResetEmailService struct {
	webURL             *url.URL
	digitalSignService identityinteractors.DigitalSignatureService
	logger             application.Logger
}

func NewRequestPassordResetEmailService(
	config infrastructure.ConfigurationProvider,
	digitalSignService identityinteractors.DigitalSignatureService,
	logger application.Logger) RequestPassordResetEmailService {
	weburlstr := config.Get("APP_URL", "https://gear5th.com")
	w, err := url.Parse(weburlstr)
	if err != nil {
		panic(err.Error())
	}
	return RequestPassordResetEmailService{
		w,
		digitalSignService,
		logger,
	}
}

func (s RequestPassordResetEmailService) SendMail(u user.User, passwordResetToken string) error {

	passwordResetURL := s.buildPasswordResetWebURL(u, passwordResetToken)

	//TODO send mail here

	s.logger.Info("mail.resetpassword", fmt.Sprintf("Sent reset password mail to %s <-> %s \n", u.Email().String(), passwordResetURL))

	return nil
}

func (r RequestPassordResetEmailService) buildPasswordResetWebURL(u user.User, passwordResetToken string) string {
	// <APP_URL>/publish/identity/managed/{userID}/reset-password?token={passwordResetToken}
	r.webURL.Path = fmt.Sprintf("/publish/identity/managed/%s/reset-password", u.UserID().String())
	q := r.webURL.Query()
	q.Set("token", passwordResetToken)
	r.webURL.RawQuery = q.Encode()
	return r.webURL.String()
}
