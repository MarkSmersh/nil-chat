package utils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	secret := GetJWTSecret()

	signed, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err.Error())
	}

	return signed
}

func DecodeToken(signed string) (jwt.Token, error) {
	secret := GetJWTSecret()

	token, err := jwt.Parse(signed, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		// jwt.WithExpirationRequired(),
		// jwt.WithLeeway(time.Hour*24*365),
	)

	// FIXME: very bad pointer (destroys the server)
	if err != nil {
		return *token, err
	}

	return *token, nil
}

// returns a new refresh token and it's iss (Issuer) claim
func NewRefreshToken(userId int) (string, string) {
	iss := uuid.New().String()

	return NewToken(jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userId),
		"iss": iss,
		// 1 year
		"exp": time.Now().Unix() + 60*60*24*365,
	}), iss
}

// returns iss (Issuer) and sub (Subject) claims of the token
func DecodeRefreshToken(refreshToken string) (string, string, error) {
	token, err := DecodeToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	iss, _ := token.Claims.GetIssuer()
	sub, _ := token.Claims.GetSubject()

	return iss, sub, nil
}

func NewAccessToken(refreshToken string) (string, error) {
	iss, sub, err := DecodeRefreshToken(refreshToken)

	if err != nil {
		return "", err
	}

	return NewToken(jwt.MapClaims{
		"iss": iss,
		"sub": sub,
		"exp": time.Now().Unix() + 60*60,
	}), nil
}

// returns iss (Issuer) and sub (Subject) claims of the token
func DecodeAccessToken(accessToken string) (string, int, error) {
	token, err := DecodeToken(accessToken)

	iss, _ := token.Claims.GetIssuer()
	subStr, _ := token.Claims.GetSubject()

	sub, _ := strconv.Atoi(subStr)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return iss, sub, err
		}

		return "", 0, err
	}

	return iss, sub, nil
}
