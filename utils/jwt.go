package utils

import (
	"backendUAS/app/models/postgres"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretString = []byte("string rahasia")

func CreateToken(user postgres.User) (string, error) {

	claims := postgres.JWTClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Role:     user.RoleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1 jam saja
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secretString))
	if err != nil {
		return "", err
	}

	return signed, err
}

// membuat token jwt yang expired dalam 24 jam
func RefreshToken(user postgres.User) (string, error) {

	claims := postgres.JWTClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Role:     user.RoleID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secretString))
	if err != nil {
		return "", err
	}

	return signed, err
}

func ValidateToken(tokenString string) (*postgres.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &postgres.JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretString, nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*postgres.JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
