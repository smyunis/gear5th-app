package accesstoken

import (
	"time"

	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type JwtAccessTokenGenenrator struct{}

func NewJwtAccessTokenGenenrator() JwtAccessTokenGenenrator {
	return JwtAccessTokenGenenrator{}
}

func (j JwtAccessTokenGenenrator) Generate(subject shared.Id) (string, error) {
	appDomain := env.Get("APP_DOMAIN", "api.gear5th.com")
	tokenClaims := jwt.RegisteredClaims{
		Subject:   subject.String(),
		Issuer:    appDomain,
		Audience:  jwt.ClaimStrings{appDomain},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	keyStr := env.Get("ACCESS_TOKEN_SIGNING_KEY", "jlPoFFmLUpLRf44w1vU9mJSXvRJleg2gk1bEzaaJN1cpafWhjyU0K4D1sOek3gDxfHWLRKrJ")
	key := []byte(keyStr)

	accessToken, err := token.SignedString(key)

	if err != nil {
		return "", err
	}
	return accessToken, nil
}
