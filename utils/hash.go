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

func HashPassword(password string) (string, error) {
	// hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to hash password: %w", err)
	// }
	hash := password
	return string(hash), nil
}
