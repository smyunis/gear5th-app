package googleoauth

import (
	"fmt"

	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

type GoogleOAuthServiceImpl struct {
	httpClient infrastructure.HTTPClient
}

func NewGoogleOAuthService(httpClient infrastructure.HTTPClient) GoogleOAuthServiceImpl {
	return GoogleOAuthServiceImpl{
		httpClient: httpClient,
	}
}

func (g *GoogleOAuthServiceImpl) UserDetails(u user.OAuthUser) user.GoogleOAuthUserDetails {
	// TODO Fetch user details from google apis
	return user.GoogleOAuthUserDetails{}
}

func (g *GoogleOAuthServiceImpl) ValidateUser(u *user.User, identityToken string) error {
	// TODO validate user token using google apis

	if u.AuthenticationMethod() != user.OAuth {
		return fmt.Errorf("user is not signed up with google oauth")
	}

	return nil
}
