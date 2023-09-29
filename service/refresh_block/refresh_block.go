package refresh_block

import (
	"block_chain/block_structure/blockchain"
	transactionUtils "block_chain/block_structure/utils"
	"block_chain/database/block"
	"block_chain/database/block_control"
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"container/heap"
	"fmt"
	"time"
	"unsafe"
)

func RefreshBlock(BlockTransactionChannel chan blockchain.BlockTransaction, MinersBlockChannel chan blockchain.Block, CompleteBlockChannel chan blockchain.Block) {
	TransactionPool := make(utils.PriorityQueue, 0)
	var TargetBlock blockchain.Block
	key, err := utils.DecodePublicKey(utils.GetPublicKey())
	if err != nil {
		return
	}
	address := utils.GetAddress(conseous.VersionPrefix, key)
	minerTransaction, err := transactionUtils.MinerTransaction(address)
	if err != nil {
		return
	}
	for {
		select {
		case bt := <-BlockTransactionChannel:
			fmt.Println("Receive Transaction:", bt.TransactionHash, " in RB")
			heap.Push(&TransactionPool, &bt)
			var tmpBlock blockchain.Block
			tmpBlock.BlockTransactions = append(tmpBlock.BlockTransactions, minerTransaction)
			tmpTransactionPool := TransactionPool
			for tmpTransactionPool.Len() > 0 {
				t := heap.Pop(&tmpTransactionPool).(*blockchain.BlockTransaction)
				tmpBlock.BlockTransactions = append(tmpBlock.BlockTransactions, *t)
				if unsafe.Sizeof(tmpBlock) > conseous.BlockSize {
					tmpBlock.BlockTransactions = tmpBlock.BlockTransactions[:len(tmpBlock.BlockTransactions)-1]
					break
				}
			}
			tmpBlock.ComputeMarkleRoot()
			fmt.Println("Complete Markle Root:", tmpBlock.BlockTop.MarkleRoot, " in RB")
			if tmpBlock.BlockTop.MarkleRoot != TargetBlock.BlockTop.MarkleRoot {
				TargetBlock = tmpBlock
				TargetBlock.BlockTop.Version = conseous.Version
				TargetBlock.BlockTop.TimeStamp = time.Now().Unix()
				blockHash, err := block_control.GetLastBlock()

				if err != nil || blockHash == conseous.GenesisBlockPreviousHash {
					blockHash = conseous.GenesisBlockPreviousHash
					TargetBlock.BlockMeta = blockchain.BlockMeta{Content: conseous.GenesisBlockContent}
					TargetBlock.BlockTop.BlockHeight = conseous.GenesisBlockHeight
					TargetBlock.BlockTop.Difficulty = conseous.GenesisBlockDifficulty
				} else {
					b, err := block.GetBlock(blockHash)
					fmt.Println(b.BlockHash, err)
					if err != nil {
						continue
					}
					TargetBlock.BlockTop.BlockHeight = b.BlockTop.BlockHeight + 1
					TargetBlock.BlockTop.Difficulty = b.BlockTop.Difficulty
					if TargetBlock.BlockTop.BlockHeight%conseous.DifficultyCycle == 0 {
						speed, err := block.CheckGenerateSpeed(TargetBlock.BlockTop.BlockHeight)
						if err == nil {
							if speed {
								if TargetBlock.BlockTop.Difficulty > conseous.DifficultyLower {
									TargetBlock.BlockTop.Difficulty -= 1
								}
							} else {
								if TargetBlock.BlockTop.Difficulty < conseous.DifficultyUpper {
									TargetBlock.BlockTop.Difficulty += 1
								}
							}
						}
					}
				}
				TargetBlock.BlockTop.PreviousHash = blockHash
				fmt.Println("Sent Block in RB")
				fmt.Println()
				MinersBlockChannel <- TargetBlock
			}
		case cb := <-CompleteBlockChannel:
			fmt.Println("Receive Complete Block", cb.BlockHash, " in RB")
			TargetBlock = blockchain.Block{}
			TargetBlock.BlockTop.Version = conseous.Version
			TargetBlock.BlockTop.TimeStamp = time.Now().Unix()
			TargetBlock.BlockTop.PreviousHash = cb.BlockHash
			TargetBlock.BlockTop.BlockHeight = cb.BlockTop.BlockHeight + 1
			TargetBlock.BlockTop.Difficulty = cb.BlockTop.Difficulty
			if TargetBlock.BlockTop.BlockHeight%conseous.DifficultyCycle == 0 {
				speed, err := block.CheckGenerateSpeed(TargetBlock.BlockTop.BlockHeight)
				if err == nil {
					if speed {
						if TargetBlock.BlockTop.Difficulty > conseous.DifficultyLower {
							TargetBlock.BlockTop.Difficulty -= 1
						}
					} else {
						if TargetBlock.BlockTop.Difficulty < conseous.DifficultyUpper {
							TargetBlock.BlockTop.Difficulty += 1
						}
					}
				}
			}
			TargetBlock.BlockTransactions = append(TargetBlock.BlockTransactions, minerTransaction)
			tmpTransactionPool := make(utils.PriorityQueue, 0)
			for TransactionPool.Len() > 0 {
				t := heap.Pop(&TransactionPool).(*blockchain.BlockTransaction)
				flag := false
				for _, transaction := range cb.BlockTransactions {
					if transaction.TransactionHash == t.TransactionHash {
						flag = true
					}
				}
				if !flag {
					heap.Push(&tmpTransactionPool, t)
					TargetBlock.BlockTransactions = append(TargetBlock.BlockTransactions, *t)
					if unsafe.Sizeof(TargetBlock) > conseous.BlockSize {
						TargetBlock.BlockTransactions = TargetBlock.BlockTransactions[:len(TargetBlock.BlockTransactions)-1]
					}
				}
			}
			TargetBlock.ComputeMarkleRoot()
			fmt.Println("New Block Markle root", TargetBlock.BlockTop.MarkleRoot, " in RB")
			TransactionPool = tmpTransactionPool
			fmt.Println("Sent New Block in RB")
			fmt.Println()
			MinersBlockChannel <- TargetBlock
		default:
			continue
		}
	}
}
