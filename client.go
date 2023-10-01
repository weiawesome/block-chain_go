package main

import (
	"block_chain/api"
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"fmt"
)

func main() {
	pair, err := utils.GenerateKeyPair()
	if err != nil {
		return
	}
	fmt.Println("PublicKey: ", pair.PublicKey)
	fmt.Println("PrivateKey: ", pair.PrivateKey)
	key, err := utils.DecodePublicKey(pair.PublicKey)
	if err != nil {
		return
	}
	addr := utils.GetAddress(conseous.VersionPrefix, key)
	fmt.Println("Address: ", addr)

	clientAPI := api.ClientAPI{ConnectNodeAddr: "127.0.0.1:8081"}
	err = clientAPI.Connect()
	if err != nil {
		return
	}
	defer func() {
		err := clientAPI.DisConnect()
		if err != nil {
			return
		}
	}()

	err = clientAPI.SubmitFreeTransaction(100, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}
