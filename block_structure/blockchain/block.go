package blockchain

import "block_chain/block_structure/blockchain/utils"

type Block struct {
	BlockTop          BlockTop           `json:"blockTop"`
	BlockTransactions []BlockTransaction `json:"blockTransactions"`
	BlockHash         string             `json:"blockHash"`
	BlockMeta         BlockMeta          `json:"blockMeta"`
}

func (b *Block) ComputeMarkleRoot() {
	var HashList []string
	for _, transaction := range b.BlockTransactions {
		HashList = append(HashList, transaction.TransactionHash)
	}
	for len(HashList) != 1 {
		var tmpHashList []string
		for i := range HashList {
			if i%2 == 0 {
				if i == len(HashList)-1 {
					tmpHashList = append(tmpHashList, utils.HashSHA256(HashList[i]))
				} else {
					tmpHashList = append(tmpHashList, utils.HashSHA256(HashList[i])+utils.HashSHA256(HashList[i+1]))
				}
			}
		}
		HashList = tmpHashList
	}
	b.BlockTop.MarkleRoot = HashList[0]
}

func (b *Block) ComputeBlockHash() {
	b.BlockHash = utils.HashSHA256(utils.HashSHA256(b.BlockTop.ToString()))
}

func (b *Block) CheckDifficulty() bool {
	for i := int64(0); i < b.BlockTop.Difficulty; i++ {
		if b.BlockHash[i] != '0' {
			return false
		}
	}
	return true
}
