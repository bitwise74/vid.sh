package util

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken(n int) (string, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
