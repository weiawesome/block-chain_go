package main

import (
	"block_chain/api"
)

func main() {
	var Amount float64
	Amount = 0

	var Addr string
	Addr = "Address"

	ConnectAddr := "127.0.0.1:8080"

	clientAPI := api.ClientAPI{ConnectNodeAddr: ConnectAddr}
	err := clientAPI.Connect()
	if err != nil {
		return
	}
	defer func() {
		err := clientAPI.DisConnect()
		if err != nil {
			return
		}
	}()

	err = clientAPI.SubmitFreeTransaction(Amount, Addr)
	if err != nil {
		return
	}
}
