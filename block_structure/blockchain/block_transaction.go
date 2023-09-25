package blockchain

import (
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
