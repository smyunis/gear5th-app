package user

import "gitlab.com/gear5th/gear5th-app/internal/domain/shared"

type GoogleOAuthUserDetails struct {
	Fullname  string
	AccountID string
	Email     string
}

type GoogleOAuthService interface {
	ValidateToken(identityToken string) (GoogleOAuthUserDetails, error)
}

type OAuthUserRepository interface {
	shared.EntityRepository[OAuthUser]
	UserWithAccountID(accountID string) (OAuthUser, error)
}

type OAuthProvider int

const (
	_ OAuthProvider = iota
	GoogleOAuth
	GithubOAuth
)

type OAuthUser struct {
	userID        shared.ID
	userAccountID string
	oauthProvider OAuthProvider
}

func ReconstituteOAuthUser(
	userID shared.ID,
	userAccountID string,
	oauthProvider OAuthProvider) OAuthUser {
	return OAuthUser{
		userID,
		userAccountID,
		oauthProvider,
	}
}

func (o *OAuthUser) UserID() shared.ID {
	return o.userID
}

func (o *OAuthUser) UserAccountID() string {
	return o.userAccountID
}

func (o *OAuthUser) OAuthProvider() OAuthProvider {
	return o.oauthProvider
}
