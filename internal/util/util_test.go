package util

import (
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashRefreshTokenSuccess(t *testing.T) {
	token := "valid-refresh-token-123"

	hashedToken, err := HashRefreshToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, hashedToken)
	require.NotEqual(t, token, hashedToken)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token)))
}

func TestGenerateAccessTokenSuccess(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")

	clientIP := "192.168.1.1"

	token, err := GenerateAccessToken(clientIP)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims := parsedToken.Claims.(jwt.MapClaims)
	assert.Equal(t, clientIP, claims["clientIP"])
	assert.NotEmpty(t, claims["exp"])
}

func TestGenerateRandomStringWithPositiveLength(t *testing.T) {
	length := 10

	result, err := GenerateRandomString(length)

	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Len(t, result, 16)
}
