package main

import (
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"fmt"
)

func main() {
	pair, err := utils.GenerateKeyPair()
	if err != nil {
		return
	}

	privateKey, err := utils.DecodePrivateKey(pair.PrivateKey)
	if err != nil {
		return
	}
	key, err := utils.DecodePublicKey(pair.PublicKey)
	if err != nil {
		return
	}

	publicKey := key

	bitcoinAddress := utils.GetAddress(conseous.VersionPrefix, publicKey)
	fmt.Println("私钥:", privateKey.D)

	fmt.Println("比特币地址:", bitcoinAddress)
}
