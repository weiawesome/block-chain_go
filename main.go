package main

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"block_chain/protocal/conseous"
	"block_chain/service/miners_leader"
	"block_chain/service/receive_validate_block"
	"block_chain/service/receive_validate_transaction"
	"block_chain/service/refresh_block"
	"block_chain/service/refresh_db"
	"block_chain/utils"
	"fmt"
)

func main() {
	fmt.Println()
	err := utils.InitClient("localhost:27017")
	if err != nil {
		return
	}
	TransactionChannel := make(chan transaction.Transaction)
	BroadcastTransactionChannel := make(chan transaction.Transaction)
	BlockTransactionChannel := make(chan blockchain.BlockTransaction)
	BroadcastBlockChannel := make(chan blockchain.Block)
	MinersSuccessBlockChannel := make(chan blockchain.Block)
	BlockChannel := make(chan blockchain.Block)
	RefreshBlockChannel := make(chan blockchain.Block)
	MinersBlockChannel := make(chan blockchain.Block)
	CompleteBlockChannel := make(chan blockchain.Block)
	go receive_validate_transaction.ReceiveValidateTransaction(TransactionChannel, BroadcastTransactionChannel, BlockTransactionChannel)
	go receive_validate_block.ReceiveValidateBlock(BroadcastBlockChannel, MinersSuccessBlockChannel, BlockChannel, RefreshBlockChannel)
	go refresh_block.RefreshBlock(BlockTransactionChannel, MinersBlockChannel, CompleteBlockChannel)
	go refresh_db.RefreshDb(RefreshBlockChannel, CompleteBlockChannel)
	go miners_leader.MinerLeader(1, MinersBlockChannel, MinersSuccessBlockChannel)
	key := utils.GetPublicKey()
	publicKey, err := utils.DecodePublicKey(key)
	if err != nil {
		return
	}
	address := utils.GetAddress(conseous.VersionPrefix, publicKey)
	var to []transaction.To
	to = append(to, transaction.To{Address: address, Amount: float64(10)})
	var from []transaction.From
	from = append(from, transaction.From{UTXOHash: conseous.TestForUXTOHash})
	t := transaction.Transaction{To: to, Fee: float64(50), PublicKey: key, From: from}
	val, err := t.ToString()
	if err != nil {
		return
	}
	k := utils.GetPrivateKey()
	privateKey, err := utils.DecodePrivateKey(k)
	if err != nil {
		return
	}
	t.TransactionHash = utils.HashSHA256(val)
	signature, err := utils.Signature(val, privateKey)
	if err != nil {
		return
	}
	t.Signature = signature
	TransactionChannel <- t
	select {}
}
