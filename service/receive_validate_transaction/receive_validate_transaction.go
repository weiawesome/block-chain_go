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
				transactionValue, err := t.ToString()
				if err != nil {
					continue
				}

				SignatureStatus := false

				var address string
				if t.PublicKey == conseous.MasterPublicKey && t.Signature == conseous.MasterSignature {
					fmt.Println("Transaction is master in VT")
					address = conseous.MasterAddress
					SignatureStatus = true
				} else {
					fmt.Println("Transaction is not master in VT")
					key, err := utils.DecodePublicKey(t.PublicKey)
					if err != nil {
						continue
					}
					address = utils.GetAddress(conseous.VersionPrefix, key)
					if utils.Verify(t.Signature, key, transactionValue) {
						fmt.Println("Transaction Verify passed in VT")
						SignatureStatus = true
					}
				}

				if SignatureStatus {
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
						fmt.Println("Fail due to not owner in VT")
						continue
					}
					sumOfSpent := float64(0)
					for _, to := range t.To {
						sumOfSpent += to.Amount
					}
					sumOfSpent += t.Fee
					fmt.Println("User have: ", sumOfHave, " in VT")
					fmt.Println("User spent: ", sumOfSpent, " in VT")
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
