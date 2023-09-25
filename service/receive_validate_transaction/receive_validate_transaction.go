package receive_validate_transaction

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	TransactionUtils "block_chain/block_structure/utils"
	blockdb "block_chain/database/block"
	"block_chain/database/block_control"
	"block_chain/database/utxo"
	"block_chain/protocal/conseous"
	"block_chain/utils"
)

func ReceiveTransaction(TransactionChannel chan transaction.Transaction, BroadcastTransactionChannel chan transaction.Transaction, BlockTransactionChannel chan blockchain.BlockTransaction) {
	for {
		select {
		case t := <-TransactionChannel:
			block, err := block_control.GetLastBlock()
			if err != nil {
				continue
			}
			if blockdb.TransactionInCheckedBlocks(block, t.TransactionHash) {
				StringValue, err := t.ToString()
				if err != nil {
					continue
				}
				key, err := utils.DecodePublicKey(t.PublicKey)
				if err != nil {
					continue
				}
				address := utils.GetAddress(conseous.VersionPrefix, key)
				if utils.Verify(t.Signature, key, StringValue) {
					sumOfHave := float64(0)
					flag := false
					for _, from := range t.From {
						value, err := utxo.GetUTXO(from.UTXOHash, from.Index)
						if err != nil || value.Spent == true {
							continue
						}
						if address != value.Address {
							flag = true
							break
						}
						sumOfHave += value.Amount
					}
					if flag {
						continue
					}
					sumOfSpent := float64(0)
					for _, to := range t.To {
						sumOfSpent += to.Amount
					}
					sumOfSpent += t.Fee

					if sumOfSpent > sumOfHave {
						continue
					}
					BroadcastTransactionChannel <- t
					BlockTransactionChannel <- TransactionUtils.ConvertTransaction(t)
				}
			}
		default:
			continue
		}
	}
}
