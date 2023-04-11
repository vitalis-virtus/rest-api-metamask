package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vitalis-virtus/rest-api-metamask/utils/fail"
)

type JwtHmacProvider struct {
	hmacSecret []byte
	issuer     string
	duration   time.Duration
}

func NewJwtHmacProvider(hmacSecret string, issuer string, duration time.Duration) *JwtHmacProvider {
	return &JwtHmacProvider{
		hmacSecret: []byte(hmacSecret),
		issuer:     issuer,
		duration:   duration,
	}
}

func (j *JwtHmacProvider) CreateStandard(subject string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.duration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.hmacSecret)
}

func (j *JwtHmacProvider) Verify(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return j.hmacSecret, nil
	})
	if err != nil {
		return nil, fail.ErrAuthError
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fail.ErrAuthError
}
