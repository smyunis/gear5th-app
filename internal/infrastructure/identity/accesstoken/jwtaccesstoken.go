package accesstoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-api/internal/infrastructure"
)

type JwtAccessTokenGenerator struct {
	config infrastructure.ConfigurationProvider
}

func NewJwtAccessTokenGenenrator(config infrastructure.ConfigurationProvider) JwtAccessTokenGenerator {
	return JwtAccessTokenGenerator{
		config,
	}
}

func (j JwtAccessTokenGenerator) Generate(subject shared.ID) (string, error) {
	appDomain := j.config.Get("APP_URL", "api.gear5th.com")
	tokenClaims := jwt.RegisteredClaims{
		Subject:   subject.String(),
		Issuer:    appDomain,
		Audience:  jwt.ClaimStrings{appDomain},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	keyStr := j.config.Get("ACCESS_TOKEN_SIGNING_KEY", "jlPoFFmLUpLRf44w1vU9mJSXvRJleg2gk1bEzaaJN1cpafWhjyU0K4D1sOek3gDxfHWLRKrJ")
	key := []byte(keyStr)

	accessToken, err := token.SignedString(key)

	if err != nil {
		return "", err
	}
	return accessToken, nil
}
