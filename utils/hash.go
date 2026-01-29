package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
)

func HashPassword(plainPassword string) string {
	secret := GetHashSecret()

	mac := hmac.New(sha256.New, []byte(secret))

	_, err := mac.Write([]byte(plainPassword))

	if err != nil {
		log.Fatal(err.Error())
	}

	hashedPassword := hex.EncodeToString(
		mac.Sum(nil),
	)

	return hashedPassword
}
