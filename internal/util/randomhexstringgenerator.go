package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func GenerateRandomHexString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be greater than 0")
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil

}
