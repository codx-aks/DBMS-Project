package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.New("failed to generate session token")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
