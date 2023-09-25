package refresh_block

import (
	"block_chain/block_structure/blockchain"
	transactionUtils "block_chain/block_structure/utils"
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"container/heap"
	"fmt"
	"unsafe"
)

func RefreshBlock(BlockTransactionChannel chan blockchain.BlockTransaction, BlockChannel chan blockchain.Block, CompleteBlockChannel chan blockchain.Block) {
	TransactionPool := make(utils.PriorityQueue, 0)
	for {
		select {
		case bt := <-BlockTransactionChannel:
			heap.Push(&TransactionPool, bt)
			var TargetBlock blockchain.Block
			key, err := utils.DecodePublicKey(utils.GetPublicKey())
			if err != nil {
				return
			}
			address := utils.GetAddress(conseous.VersionPrefix, key)
			transaction, err := transactionUtils.MinerTransaction(address)
			if err != nil {
				continue
			}
			TargetBlock.BlockTransactions = append(TargetBlock.BlockTransactions, transaction)
			tmpTransactionPool := TransactionPool
			for tmpTransactionPool.Len() > 0 {
				t := heap.Pop(&tmpTransactionPool).(*blockchain.BlockTransaction)
				TargetBlock.BlockTransactions = append(TargetBlock.BlockTransactions, *t)
				if unsafe.Sizeof(TargetBlock) > conseous.BlockSize {
					TargetBlock.BlockTransactions = TargetBlock.BlockTransactions[:len(TargetBlock.BlockTransactions)-1]
					break
				}
			}
			TargetBlock.ComputeMarkleRoot()
		case cb := <-CompleteBlockChannel:
			fmt.Println(cb)
		default:
			continue
		}
	}
}
