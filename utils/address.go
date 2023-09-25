package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

func GetAddress(networkPrefix byte, publicKey *ecdsa.PublicKey) string {
	publicKeyHash := publicKeyHash(publicKey)
	payload := append([]byte{networkPrefix}, publicKeyHash...)
	checksum := checksum(payload)
	payload = append(payload, checksum...)
	address := base58Encode(payload)
	return address
}

func publicKeyHash(publicKey *ecdsa.PublicKey) []byte {
	pubKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
	sha256Hash := sha256.Sum256(pubKeyBytes)
	h := ripemd160.New()
	h.Write(sha256Hash[:])
	publicKeyHash := h.Sum(nil)
	return publicKeyHash
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:4]
}

func base58Encode(input []byte) string {
	base58Alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	var result []byte

	x := new(big.Int).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)

	for x.Cmp(zero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, base, mod)
		result = append([]byte{base58Alphabet[mod.Int64()]}, result...)
	}

	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{base58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return string(result)
}
