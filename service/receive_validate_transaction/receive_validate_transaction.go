package receive_validate_transaction

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	transactionUtils "block_chain/block_structure/utils"
	blockdb "block_chain/database/block"
	"block_chain/database/block_control"
	"block_chain/database/utxo"
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"fmt"
)

func ReceiveValidateTransaction(TransactionChannel chan transaction.Transaction, BroadcastTransactionChannel chan transaction.Transaction, BlockTransactionChannel chan blockchain.BlockTransaction) {
	for {
		select {
		case t := <-TransactionChannel:
			fmt.Println("Receive Transaction:", t.TransactionHash, " in VT")
			block, err := block_control.GetLastBlock()
			if err != nil {
				block = conseous.GenesisBlockPreviousHash
			}
			fmt.Println("Last block hash:", block, " in VT")
			if blockdb.TransactionInCheckedBlocks(block, t.TransactionHash) {
				fmt.Println("Pass Transaction checked in block in VT")
				StringValue, err := t.ToString()
				if err != nil {
					continue
				}
				key, err := utils.DecodePublicKey(t.PublicKey)
				if err != nil {
					continue
				}
				address := utils.GetAddress(conseous.VersionPrefix, key)
				fmt.Println("Verify signature:", utils.Verify(t.Signature, key, StringValue), " in VT")
				if utils.Verify(t.Signature, key, StringValue) || t.Signature == conseous.MasterSignature || t.PublicKey == conseous.MasterPublicKey {
					sumOfHave := float64(0)
					flag := false
					for _, from := range t.From {
						value, err := utxo.GetUTXO(from.UTXOHash, from.Index)
						if err != nil || value.Spent == true {
							continue
						}
						if address != value.Address && value.Address != conseous.MasterAddress {
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
					fmt.Println("Have and spent:", sumOfHave, sumOfSpent)
					if sumOfSpent > sumOfHave {
						continue
					}
					fmt.Println("Sent Transaction:", t.TransactionHash, " in VT")
					fmt.Println()
					BroadcastTransactionChannel <- t
					BlockTransactionChannel <- transactionUtils.ConvertTransaction(t)
				}
			}
		default:
			continue
		}
	}
}
