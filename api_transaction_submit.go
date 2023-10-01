package main

import (
	"block_chain/api"
	"block_chain/block_structure/transaction"
	"block_chain/utils"
	"fmt"
)

func main() {
	var UTXOHash string
	var Index int
	var Address string
	var Amount float64
	var Fee float64
	var PublicKey string
	var PrivateKey string

	UTXOHash = "UTXOHash"
	Index = 0
	Address = "Address"
	Amount = 0
	Fee = 0
	PublicKey = "PublicKey"
	PrivateKey = "PrivateKey"

	ConnectAddr := "127.0.0.1:8080"

	f := transaction.From{UTXOHash: UTXOHash, Index: Index}
	t := transaction.To{Address: Address, Amount: Amount}

	ts := transaction.Transaction{From: []transaction.From{f}, To: []transaction.To{t}, Fee: Fee, PublicKey: PublicKey}

	privateKey, err := utils.DecodePrivateKey(PrivateKey)
	if err != nil {
		return
	}
	strVal, err := ts.ToString()
	if err != nil {
		return
	}
	signature, err := utils.Signature(strVal, privateKey)
	if err != nil {
		return
	}
	ts.Signature = signature
	ts.TransactionHash = utils.HashSHA256(strVal)
	fmt.Println("TransactionHash: ", ts.TransactionHash)
	fmt.Println()
	fmt.Println("From: ")
	fmt.Println("PublicKey: ", PublicKey)
	fmt.Println("PrivateKey: ", PrivateKey)
	fmt.Println("UTXOHash: ", UTXOHash)
	fmt.Println("Index: ", Index)
	fmt.Println()
	fmt.Println("To: ")
	fmt.Println("Address: ", Address)
	fmt.Println("Amount: ", Amount)
	fmt.Println()
	fmt.Println("Fee: ", Fee)
	fmt.Println()
	fmt.Println("Signature: ", signature)

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

	err = clientAPI.SubmitTransaction(ts.From, ts.To, ts.Fee, ts.PublicKey, ts.Signature)
	if err != nil {
		fmt.Println(err)
		return
	}
}
