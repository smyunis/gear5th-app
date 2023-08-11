package identityemail

import (
	"fmt"
	"net/url"

	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
)

type RequestPassordResetEmailService struct {
	webURL *url.URL
}

func NewRequestPassordResetEmailService(config infrastructure.ConfigurationProvider) RequestPassordResetEmailService {
	weburlstr := config.Get("PUBLISHER_WEB_URL", "https://publisher.gear5th.com")
	a, err := url.Parse(weburlstr)
	if err != nil {
		panic(err.Error())
	}
	return RequestPassordResetEmailService{
		webURL: a,
	}
}

func (r RequestPassordResetEmailService) SendMail(u user.User, passwordResetToken string) error {

	//https://publisher.gear5th.com/identity/managed/{userID}/reset-password?reset-token={passwordResetToken}

	r.webURL.Path = fmt.Sprintf("/identity/managed/%s/reset-password", u.UserID().String())
	q := r.webURL.Query()
	q.Set("token", passwordResetToken)
	r.webURL.RawQuery = q.Encode()

	//TODO send mail here
	fmt.Printf("Sent mail to %s <-> %s \n", u.Email().String(), r.webURL.String())

	return nil
}
