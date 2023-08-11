package identityemail

import (
	"fmt"
	"net/url"
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
)

type RequestPassordResetEmailService struct {
	webURL  *url.URL
	kvStore application.KeyValueStore
}

func NewRequestPassordResetEmailService(config infrastructure.ConfigurationProvider,
	kvStore application.KeyValueStore) RequestPassordResetEmailService {
	weburlstr := config.Get("PUBLISHER_WEB_URL", "https://publisher.gear5th.com")
	w, err := url.Parse(weburlstr)
	if err != nil {
		panic(err.Error())
	}
	return RequestPassordResetEmailService{
		w,
		kvStore,
	}
}

func (r RequestPassordResetEmailService) SendMail(u user.User, passwordResetToken string) error {

	passwordResetURL := r.buildPasswordResetWebURL(u, passwordResetToken)

	k := manageduserinteractors.PasswordResetTokenStoreKey(u.UserID().String())
	err := r.kvStore.Save(k, passwordResetToken, 30*time.Minute)
	if err != nil {
		return err
	}

	//TODO send mail here
	fmt.Printf("Sent reset password mail to %s <-> %s \n", u.Email().String(), passwordResetURL)

	return nil
}

func (r RequestPassordResetEmailService) buildPasswordResetWebURL(u user.User, passwordResetToken string) string {
	// <PUBLISHER_WEB_URL>/identity/managed/{userID}/reset-password?token={passwordResetToken}
	r.webURL.Path = fmt.Sprintf("/identity/managed/%s/reset-password", u.UserID().String())
	q := r.webURL.Query()
	q.Set("token", passwordResetToken)
	r.webURL.RawQuery = q.Encode()
	return r.webURL.String()
}
