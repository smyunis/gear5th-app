package user

import "gitlab.com/gear5th/gear5th-api/internal/domain/shared"

type GoogleOAuthUserDetails struct {
}

type GoogleOAuthService interface {
	UserDetails(u OAuthUser) GoogleOAuthUserDetails
	ValidateUser(u *User, identityToken string) error
}

type OAuthUserRepository interface {
	shared.EntityRepository[OAuthUser]
}
type OAuthUser struct {
	userId         shared.Id
	userIdentifier string
	oauthProvider  OAuthProvider
}

type OAuthProvider int

const (
	_ OAuthProvider = iota
	GoogleOAuth
	GithubOAuth
)

func (o *OAuthUser) UserIdentifier() string {
	return o.userIdentifier
}

func (o *OAuthUser) OAuthProvider() OAuthProvider {
	return o.oauthProvider
}
