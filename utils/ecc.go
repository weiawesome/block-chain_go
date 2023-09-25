package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
)

type SignatureData struct {
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

type KeyPair struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func GenerateKeyPair() (KeyPair, error) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return KeyPair{}, err
	}
	privateKeyBytes := privateKey.D.Bytes()
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyBytes)

	publicKey := &privateKey.PublicKey
	publicKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)

	return KeyPair{PublicKey: publicKeyBase64, PrivateKey: privateKeyBase64}, err
}

func DecodePublicKey(PublicKey string) (*ecdsa.PublicKey, error) {
	curve := elliptic.P256()

	decodedPublicKeyBytes, err := base64.StdEncoding.DecodeString(PublicKey)
	if err != nil {
		return &ecdsa.PublicKey{}, err
	}

	decodedPublicKey := new(ecdsa.PublicKey)
	decodedPublicKey.Curve = curve
	decodedPublicKey.X, decodedPublicKey.Y = elliptic.Unmarshal(curve, decodedPublicKeyBytes)

	return decodedPublicKey, nil
}

func DecodePrivateKey(PrivateKey string) (*ecdsa.PrivateKey, error) {
	curve := elliptic.P256()
	decodedPrivateKeyBytes, err := base64.StdEncoding.DecodeString(PrivateKey)
	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}

	decodedPrivateKey := new(ecdsa.PrivateKey)
	decodedPrivateKey.Curve = curve
	decodedPrivateKey.D = new(big.Int).SetBytes(decodedPrivateKeyBytes)

	return decodedPrivateKey, nil
}

func Signature(Content string, PrivateKey *ecdsa.PrivateKey) (string, error) {
	ContentBytes := []byte(Content)
	r, s, err := ecdsa.Sign(rand.Reader, PrivateKey, ContentBytes)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(SignatureData{R: r, S: s})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func Verify(Signature string, PublicKey *ecdsa.PublicKey, Content string) bool {
	ContentBytes := []byte(Content)
	var signature SignatureData

	err := json.Unmarshal([]byte(Signature), &signature)
	if err != nil {
		return false
	}

	valid := ecdsa.Verify(PublicKey, ContentBytes, signature.R, signature.S)

	return valid
}
