package util

import "golang.org/x/crypto/bcrypt"

// HashRefreshToken generates a bcrypt hash of the provided refresh token.
func HashRefreshToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
