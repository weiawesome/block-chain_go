package initalize

import (
	"block_chain/block_structure/blockchain"
	"block_chain/database/utxo"
)

func BuildUTXO(Block blockchain.Block) error {
	for i, transaction := range Block.BlockTransactions {
		if i == 0 {
			continue
		}
		err := utxo.SetUTXO(transaction)
		if err != nil {
			return err
		}
	}
	return nil
}
