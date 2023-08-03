package user

type OAuthUser struct {
	*User
	userIdentifier any
	oauthProvider OAuthProvider
}

type OAuthProvider int
const (
	GoogleOAuth OAuthProvider = iota
	GithubOAuth
)

func (o *OAuthUser) UserIdentifier() any {
	return o.userIdentifier
}

func (o *OAuthUser) OAuthProvider() OAuthProvider {
	return o.oauthProvider
}