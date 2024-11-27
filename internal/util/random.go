package util

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// GenerateRandomString generate random string by input length.
func GenerateRandomString(length int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes := make([]byte, length)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		bytes[i] = letters[num.Int64()]
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
