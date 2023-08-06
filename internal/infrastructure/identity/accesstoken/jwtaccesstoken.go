package accesstoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type JwtAccessTokenGenenrator struct{}

func (j JwtAccessTokenGenenrator) Generate(subject shared.Id) (string, error) {
	tokenClaims := jwt.RegisteredClaims{
		Subject:   subject.String(),
		Issuer:    "api.gear5th.com",
		Audience:  jwt.ClaimStrings{"api.gear5th.com"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	//TODO fetch key from an external source
	key := []byte("jlPoFFmLUpLRf44w1vU9mJSXvRJleg2gk1bEzaaJN1cpafWhjyU0K4D1sOek3gDxfHWLRKrJ")

	accessToken, err := token.SignedString(key)

	if err != nil {
		return "", err
	}
	return accessToken, nil
}
