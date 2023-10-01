package main

import (
	"block_chain/api"
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"fmt"
)

func main() {

	var Amount float64

	Amount = 0

	pair, err := utils.GenerateKeyPair()
	if err != nil {
		return
	}
	key, err := utils.DecodePublicKey(pair.PublicKey)
	if err != nil {
		return
	}
	addr := utils.GetAddress(conseous.VersionPrefix, key)

	fmt.Println("PublicKey: ", pair.PublicKey)
	fmt.Println("PrivateKey: ", pair.PrivateKey)
	fmt.Println("Address: ", addr)
	fmt.Println("Amount: ", Amount)

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

	err = clientAPI.SubmitFreeTransaction(Amount, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}
