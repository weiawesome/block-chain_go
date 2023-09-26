package main

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"block_chain/protocal/connection"
	"block_chain/utils"
	"fmt"
)

func main() {
	err := utils.InitClient("127.0.0.1:27017")
	if err != nil {
		fmt.Println(err)
		return
	}
	TransactionChannel := make(chan transaction.Transaction)
	BlockChannel := make(chan blockchain.Block)
	BroadcastTransactionChannel := make(chan transaction.Transaction)
	BroadcastBlockChannel := make(chan blockchain.Block)
	var StrList []string
	StrList = append(StrList, "127.0.0.1:8080")
	go connection.BuildNode("127.0.0.1:8081", StrList, TransactionChannel, BlockChannel, BroadcastTransactionChannel, BroadcastBlockChannel)
}
