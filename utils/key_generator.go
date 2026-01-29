package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateApiAndSecretKey(MerchantID uint) (string, string, error) {
	apiRand, err := GenerateRandomString(20)
	if err != nil {
		return "", "", err
	}
	secretRand, err := GenerateRandomString(40)
	if err != nil {
		return "", "", err
	}

	apiKey := fmt.Sprintf("%d_%s", MerchantID, apiRand)
	secretKey := fmt.Sprintf("%d_%s", MerchantID, secretRand)

	return apiKey, secretKey, nil
}
