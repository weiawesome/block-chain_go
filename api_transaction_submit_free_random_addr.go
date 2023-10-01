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

	ConnectAddr := "127.0.0.1:8080"

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

	clientAPI := api.ClientAPI{ConnectNodeAddr: ConnectAddr}
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
		return
	}

}
