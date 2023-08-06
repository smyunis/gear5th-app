package identityusecases

import "github.com/golang-jwt/jwt/v5"

type IdentityTokenClaims struct {
	jwt.RegisteredClaims
	Roles  []string `json:"roles"`
}


func (i *IdentityTokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return i.RegisteredClaims.GetExpirationTime()
}

func (i *IdentityTokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return i.RegisteredClaims.GetIssuedAt()
}

func (i *IdentityTokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return i.RegisteredClaims.GetNotBefore()
}

func (i *IdentityTokenClaims) GetIssuer() (string, error) {
	return i.RegisteredClaims.GetIssuer()
}

func (i *IdentityTokenClaims) GetSubject() (string, error) {
	return i.RegisteredClaims.GetSubject()
}

func (i *IdentityTokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return i.RegisteredClaims.GetAudience()
}
