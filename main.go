package main

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"block_chain/protocal/connection"
	"block_chain/service/miners_leader"
	"block_chain/service/receive_validate_block"
	"block_chain/service/receive_validate_transaction"
	"block_chain/service/refresh_block"
	"block_chain/service/refresh_db"
	"block_chain/utils"
)

func main() {
	err := utils.InitClient("localhost:27017")
	if err != nil {
		return
	}

	Miners := 1

	NodeAddr := "127.0.0.1:8081"
	NodeAddresses := []string{}
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
	go miners_leader.MinerLeader(Miners, MinersBlockChannel, MinersSuccessBlockChannel)

	connection.BuildNode(NodeAddr, NodeAddresses, TransactionChannel, BlockChannel, BroadcastTransactionChannel, BroadcastBlockChannel)
}
