package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(value string) string {
	data := []byte(value)

	hash := sha256.New()
	hash.Write(data)
	hashedData := hash.Sum(nil)
	hashString := hex.EncodeToString(hashedData)

	return hashString
}
