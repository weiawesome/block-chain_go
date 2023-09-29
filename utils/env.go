package utils

import "os"

func GetPublicKey() string {
	PublicKey := os.Getenv("PUBLIC_KEY")
	if PublicKey == "" {
		keyPair, err := GenerateKeyPair()
		if err != nil {
			return ""
		}
		err = os.Setenv("PUBLIC_KEY", keyPair.PublicKey)
		if err != nil {
			return ""
		}
		err = os.Setenv("PRIVATE_KEY", keyPair.PrivateKey)
		if err != nil {
			return ""
		}
		PublicKey = keyPair.PublicKey
	}
	return PublicKey
}
func GetPrivateKey() string {
	PrivateKey := os.Getenv("PRIVATE_KEY")
	if PrivateKey == "" {
		keyPair, err := GenerateKeyPair()
		if err != nil {
			return ""
		}
		err = os.Setenv("PUBLIC_KEY", keyPair.PublicKey)
		if err != nil {
			return ""
		}
		err = os.Setenv("PRIVATE_KEY", keyPair.PrivateKey)
		if err != nil {
			return ""
		}
		PrivateKey = keyPair.PublicKey
	}
	return PrivateKey
}
func GetNumOfMiners() string {
	PublicKey := os.Getenv("NUM_OF_MINERS")
	if PublicKey == "" {

	}
	return PublicKey
}
