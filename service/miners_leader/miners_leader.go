package miners_leader

import (
	"block_chain/block_structure/blockchain"
	"block_chain/service/miners"
	"fmt"
)

func MinerLeader(Miners int, MinersBlockChannel chan blockchain.Block, MinersSuccessBlockChannel chan blockchain.Block) {
	var minersChannels []chan blockchain.Block
	SuccessBlockChannel := make(chan blockchain.Block)
	Lower := uint32(0)
	Upper := ^uint32(0)
	nonceRange := (Upper - Lower + 1) / uint32(Miners)
	for i := 0; i < Miners; i++ {
		minerChannel := make(chan blockchain.Block)
		go miners.Miner(Lower, Lower+uint32(i)*nonceRange, minerChannel, SuccessBlockChannel)
		minersChannels = append(minersChannels, minerChannel)
	}
	for {
		select {
		case mb := <-MinersBlockChannel:
			fmt.Println("Receive block in ML")
			fmt.Println("Sent block in ML")
			fmt.Println()
			for _, channel := range minersChannels {
				channel <- mb
			}
		case sb := <-SuccessBlockChannel:
			fmt.Println("Receive Complete Block", sb.BlockHash, " in ML")
			fmt.Println("Sent Complete Block in ML")
			fmt.Println()
			MinersSuccessBlockChannel <- sb
		}
	}
}
