package auth

import (
	"time"

	"github.com/098765432m/config"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(config.AppData.JWT.SecretKey)

type UserClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(username, email, role string) (string, error) {
	claims := UserClaims {
		Username: username,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}