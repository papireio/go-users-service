package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GetToken() (string, error) {
	bytes := make([]byte, 128)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), err
}
