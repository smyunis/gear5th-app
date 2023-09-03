package googleoauth

import (
	"context"

	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
	"google.golang.org/api/idtoken"
)

type GoogleOAuthServiceImpl struct {
}

func NewGoogleOAuthService() GoogleOAuthServiceImpl {
	return GoogleOAuthServiceImpl{
	}
}

func (g GoogleOAuthServiceImpl) ValidateToken(identityToken string) (user.GoogleOAuthUserDetails, error) {
	claims, err := idtoken.Validate(context.Background(), identityToken, "")
	if err != nil {
		return user.GoogleOAuthUserDetails{}, err
	}

	details := user.GoogleOAuthUserDetails{
		Fullname: claims.Claims["name"].(string),
		AccountID: claims.Subject,
		Email: claims.Claims["email"].(string),
	}

	return details, nil
}
