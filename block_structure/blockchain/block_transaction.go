package blockchain

import (
	"block_chain/block_structure/transaction"
	"encoding/json"
	"strconv"
)

type BlockTransaction struct {
	From            []From  `json:"from"`
	To              []To    `json:"to"`
	Fee             float64 `json:"fee"`
	TransactionHash string  `json:"transactionHash"`
}

type From struct {
	UTXOHash string
	Index    int
}

type To struct {
	Address string
	Amount  float64
}

func (bt *BlockTransaction) ToString() (string, error) {
	fromJson, err := json.Marshal(bt.From)
	if err != nil {
		return "", err
	}
	toJson, err := json.Marshal(bt.To)
	if err != nil {
		return "", err
	}
	return string(fromJson) + string(toJson) + strconv.FormatFloat(bt.Fee, 'f', -1, 64), nil
}

func ConvertTransaction(transaction transaction.Transaction) BlockTransaction {
	var bt BlockTransaction
	bt.TransactionHash = transaction.TransactionHash
	bt.Fee = transaction.Fee
	for _, from := range transaction.From {
		bt.From = append(bt.From, From{UTXOHash: from.UTXOHash, Index: from.Index})
	}
	for _, to := range transaction.To {
		bt.To = append(bt.To, To{Address: to.Address, Amount: to.Amount})
	}
	return bt
}
