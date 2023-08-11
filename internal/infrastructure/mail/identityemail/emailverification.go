package identityemail

import (
	"fmt"
	"net/url"
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/application"
	"gitlab.com/gear5th/gear5th-api/internal/application/identityinteractors/manageduserinteractors"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
)

type VerifcationEmailSender struct {
	appURL  *url.URL
	kvStore application.KeyValueStore
}

func NewVerifcationEmailSender(config infrastructure.ConfigurationProvider, kvStore application.KeyValueStore) VerifcationEmailSender {
	appurlstr := config.Get("APP_URL", "https://api.gear5th.com")
	a, err := url.Parse(appurlstr)
	if err != nil {
		panic(err.Error())
	}

	return VerifcationEmailSender{
		a,
		kvStore,
	}
}

func (s VerifcationEmailSender) SendMail(event any) {
	signedUpUser := event.(user.UserCreatedEvent)

	token := shared.NewID().String()

	k := manageduserinteractors.EmailVerificationTokenStoreKey(signedUpUser.UserId.String())

	err := s.kvStore.Save(k, token, 30*time.Minute)
	if err != nil {
		fmt.Printf("error saving to kv store : %s", err.Error())
	}
	verificationURL := s.buildEmailVerificationURL(signedUpUser, token)

	//TODO send email with link to verify email

	fmt.Printf("Sending Verification Email to %s <-> %s\n", signedUpUser.Email.String(), verificationURL)

}

func (s VerifcationEmailSender) buildEmailVerificationURL(signedUpUser user.UserCreatedEvent, token string) string {
	// <APP_URL>/identity/managed/{userId}/verify-email?token={token}
	s.appURL.Path = fmt.Sprintf("/identity/managed/%s/verify-email", signedUpUser.UserId.String())
	q := s.appURL.Query()
	q.Set("token", token)
	s.appURL.RawQuery = q.Encode()
	return s.appURL.String()
}
