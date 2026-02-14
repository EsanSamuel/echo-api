package token

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateVerificationToken() (string, error) {
	n := make([]byte, 35)
	_, err := rand.Read(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(n), nil
}
