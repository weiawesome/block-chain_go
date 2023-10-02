package miners

import (
	"block_chain/block_structure/blockchain"
	"block_chain/protocal/conseous"
	"fmt"
)

func Miner(Lower uint32, Upper uint32, MinerBlockChannel chan blockchain.Block, SuccessBlockChannel chan blockchain.Block) {
	for {
		select {
		case b := <-MinerBlockChannel:
			fmt.Println("Miners Receive in MS")
			fmt.Println("Difficulty: ", b.BlockTop.Difficulty, " in MS")
			fmt.Println("Size of transaction: ", len(b.BlockTransactions), " in MS")
			if len(b.BlockTransactions) <= 1 && conseous.MineEmpty == false {
				fmt.Println("Stop")
				continue
			}

			i := Lower
			flag := false
			for {
				select {
				case nb := <-MinerBlockChannel:
					i = Lower - 1
					b = nb
				default:
					b.BlockTop.Nonce = i
					b.ComputeBlockHash()

					fmt.Println("Mine result ", i, b.BlockHash, b.CheckDifficulty(), " in MS")

					if b.CheckDifficulty() {
						fmt.Println("Sent Hash", b.BlockHash, " in MS")
						fmt.Println()
						SuccessBlockChannel <- b
						flag = true
						break
					}
					i++
					if i == Upper {
						flag = true
						break
					}
				}
				if flag {
					break
				}
			}
		default:
			continue
		}
	}
}
