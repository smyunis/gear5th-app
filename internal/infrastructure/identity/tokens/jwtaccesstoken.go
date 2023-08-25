package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure"
)

type JwtAccessTokenService struct {
	config infrastructure.ConfigurationProvider
}

func NewJwtAccessTokenService(config infrastructure.ConfigurationProvider) JwtAccessTokenService {
	return JwtAccessTokenService{
		config,
	}
}

func (s JwtAccessTokenService) Generate(subject shared.ID) (string, error) {
	appDomain := s.config.Get("APP_URL", "gear5th.com")
	tokenClaims := jwt.RegisteredClaims{
		Subject:   subject.String(),
		Issuer:    appDomain,
		Audience:  jwt.ClaimStrings{appDomain},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)), //30 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	keyStr := s.config.Get("ACCESS_TOKEN_SIGNING_KEY", "jlPoFFmLUpLRf44w1vU9mJSXvRJleg2gk1bEzaaJN1cpafWhjyU0K4D1sOek3gDxfHWLRKrJ")
	key := []byte(keyStr)

	accessToken, err := token.SignedString(key)

	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (s JwtAccessTokenService) Validate(token string) bool {
	t, err := s.verifiedToken(token)
	if err != nil {
		return false
	}
	return t.Valid
}

func (s JwtAccessTokenService) UserID(token string) (shared.ID, error) {
	t, err := s.verifiedToken(token)
	if err != nil {
		return shared.ID(""), err
	}
	userID, err := t.Claims.GetSubject()
	if err != nil {
		return shared.ID(""), err
	}
	return shared.ID(userID), nil
}

func (s JwtAccessTokenService) verifiedToken(token string) (*jwt.Token, error) {
	appDomain := s.config.Get("APP_URL", "gear5th.com")
	keyStr := s.config.Get("ACCESS_TOKEN_SIGNING_KEY", "jlPoFFmLUpLRf44w1vU9mJSXvRJleg2gk1bEzaaJN1cpafWhjyU0K4D1sOek3gDxfHWLRKrJ")
	t, err := jwt.Parse(token, func(tok *jwt.Token) (interface{}, error) {
		return []byte(keyStr), nil
	},
		jwt.WithIssuer(appDomain),
		jwt.WithAudience(appDomain),
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}
