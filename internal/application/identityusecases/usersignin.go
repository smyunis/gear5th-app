package identityusecases

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/gear5th/gear5th-api/internal/domain/identity/user"
)

type ManagedUserInteractor struct {
	userRepository        user.UserRepository
	managedUserRepository user.ManagedUserRepository
}

func NewManagedUserInteractor(userRepository user.UserRepository,
	managedUserRepository user.ManagedUserRepository) ManagedUserInteractor {
	return ManagedUserInteractor{
		userRepository:        userRepository,
		managedUserRepository: managedUserRepository,
	}
}

var ErrAuthorization = errors.New("authentication error")

func (m *ManagedUserInteractor) SignIn(email user.Email, password string) (string, error) {

	u, err := m.userRepository.UserWithEmail(email)
	if err != nil {
		return "", ErrAuthorization
	}

	managedUser, err := m.managedUserRepository.Get(u.UserID())
	if err != nil {
		return "", ErrAuthorization
	}

	if email != u.Email() {
		return "", ErrAuthorization
	}

	if !managedUser.IsPasswordCorrect(password) {
		return "", ErrAuthorization
	}

	tokenClaims := jwt.RegisteredClaims{
		Subject:   "stub-id-xxx",
		Issuer:    "api.gear5th.com",
		Audience:  jwt.ClaimStrings{"api.gear5th.com"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	key := []byte("secretkey")
	accessToken, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}
