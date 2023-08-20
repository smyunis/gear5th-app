package tokens_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Learning tests for jwt-go package

type JwtClaims struct {
	jwt.RegisteredClaims
	Foo string `json:"foo"`
}

func TestLearningJwtCreateEmptyToken(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)
	key := []byte("secretkey")
	token.SignedString(key)

}

func TestLearningJwtCreateTokenWithClaims(t *testing.T) {

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "me",
			Subject:   "this",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
		},
		Foo: "bar",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte("secretkey")
	token.SignedString(key)

}

func TestLearningJwtVerifyToken(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJpc3MiOiJtZSIsInN1YiI6InRoaXMifQ.hHNpHZZg9IO4L5zo4fVjOElV4r6VUvmTTtEL9YCKz0I"

	// expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJtZSIsInN1YiI6InRoaXMiLCJleHAiOjE2OTEyNzM1NTYsImZvbyI6ImJhciJ9.ds_5TKkoVufay6GC0DQ8Er0lersenCu74Op63DCCq6U"

	token, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		return []byte("secretkey"), nil
	},
		jwt.WithIssuer("me"),
		jwt.WithSubject("this"),
	)

	if !token.Valid {
		t.Fatalf("validation error: %s", err.Error())
	}

}

func TestLearningJwtGetClaimFromToken(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJpc3MiOiJtZSIsInN1YiI6InRoaXMifQ.hHNpHZZg9IO4L5zo4fVjOElV4r6VUvmTTtEL9YCKz0I"

	token, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
		return []byte("secretkey"), nil
	},
		jwt.WithIssuer("me"),
		jwt.WithSubject("this"),
	)

	if sub, _ := token.Claims.GetSubject(); sub != "this" {
		t.Fatal(err)
	}

	if token.Claims.(jwt.MapClaims)["foo"] != "bar" {
		t.FailNow()
	}

}
