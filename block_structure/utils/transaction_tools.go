package utils

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/blockchain/utils"
	"block_chain/block_structure/transaction"
)

func ConvertTransaction(transaction transaction.Transaction) blockchain.BlockTransaction {
	var bt blockchain.BlockTransaction
	bt.TransactionHash = transaction.TransactionHash
	bt.Fee = transaction.Fee
	for _, from := range transaction.From {
		bt.From = append(bt.From, blockchain.From{UTXOHash: from.UTXOHash, Index: from.Index})
	}
	for _, to := range transaction.To {
		bt.To = append(bt.To, blockchain.To{Address: to.Address, Amount: to.Amount})
	}
	return bt
}

func MinerTransaction(Address string) (blockchain.BlockTransaction, error) {
	var bt blockchain.BlockTransaction
	bt.To = append(bt.To, blockchain.To{Address: Address})
	stringValue, err := bt.ToString()
	if err != nil {
		return bt, err
	}
	bt.TransactionHash = utils.HashSHA256(stringValue)
	return bt, nil
}
