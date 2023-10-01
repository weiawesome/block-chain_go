package receive_validate_block

import (
	"block_chain/block_structure/blockchain"
	"block_chain/database/block"
	"fmt"
)

func ReceiveValidateBlock(BroadcastBlockChannel chan blockchain.Block, MinersSuccessBlockChannel chan blockchain.Block, BlockChannel chan blockchain.Block, RefreshBlockChannel chan blockchain.Block) {
	for {
		select {
		case mb := <-MinersSuccessBlockChannel:
			fmt.Println("Receive Block:", mb.BlockHash, " in VB")
			fmt.Println("Sent Block:", mb.BlockHash, " in VB")
			fmt.Println()
			BroadcastBlockChannel <- mb
			RefreshBlockChannel <- mb
		case b := <-BlockChannel:
			_, err := block.GetBlock(b.BlockHash)
			if err != nil {
				if b.CheckDifficulty() {
					tmpBlockHash := b.BlockHash
					b.ComputeBlockHash()
					if tmpBlockHash == b.BlockHash {
						tmpMarkleRoot := b.BlockTop.MarkleRoot
						b.ComputeMarkleRoot()
						if tmpMarkleRoot == b.BlockTop.MarkleRoot {
							BroadcastBlockChannel <- b
							RefreshBlockChannel <- b
						}
					}
				}
			}
		default:
			continue
		}
	}
}
