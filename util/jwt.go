package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SecretKey = "secret"

func GenerateJwt(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	return claims.SignedString([]byte(SecretKey))
}

func ParseJwt(cookie string) (string, error) {
	if cookie == "" {
		return "", errors.New("empty token")
	}

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims := token.Claims.(*jwt.StandardClaims)
	return claims.Issuer, nil
}