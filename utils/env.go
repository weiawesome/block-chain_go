package utils

import "os"

func GetPublicKey() string {
	PublicKey := os.Getenv("PUBLIC_KEY")
	if PublicKey == "" {

	}
	return PublicKey
}
