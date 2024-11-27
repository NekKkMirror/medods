package util

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateAccessToken creates a new JWT access token for a user.
func GenerateAccessToken(clientIP string) (string, error) {
	claims := jwt.MapClaims{
		"clientIP": clientIP,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
