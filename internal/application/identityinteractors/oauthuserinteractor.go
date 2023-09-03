package identityinteractors

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/identity/user"
)

type OAuthUserInteractor struct {
	oAuthUserRepository user.OAuthUserRepository
	tokenGenerator      AccessTokenService
	googleOAuthService  user.GoogleOAuthService
}

func NewOAuthUserInteractor(
	oAuthUserRepository user.OAuthUserRepository,
	tokenGenerator AccessTokenService,
	googleOAuthService user.GoogleOAuthService) OAuthUserInteractor {
	return OAuthUserInteractor{
		oAuthUserRepository,
		tokenGenerator,
		googleOAuthService,
	}
}

func (i *OAuthUserInteractor) SignIn(credential string) (string, error) {
	userDetails, err := i.googleOAuthService.ValidateToken(credential)
	if err != nil {
		return "", application.ErrAuthorization
	}

	oauthUser, err := i.oAuthUserRepository.UserWithAccountID(userDetails.AccountID)
	if err != nil {
		return "",  application.ErrAuthorization
	}

	return i.tokenGenerator.Generate(oauthUser.UserID())
}
