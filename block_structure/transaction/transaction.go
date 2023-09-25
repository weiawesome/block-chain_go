package transaction

import (
	"encoding/json"
	"strconv"
)

type Transaction struct {
	From            []From
	To              []To
	Fee             float64
	Signature       string
	TransactionHash string
	PublicKey       string
}

type From struct {
	UTXOHash string
	Index    int
}

type To struct {
	Address string
	Amount  float64
}

func (t Transaction) ToString() (string, error) {
	fromJson, err := json.Marshal(t.From)
	if err != nil {
		return "", err
	}
	toJson, err := json.Marshal(t.To)
	if err != nil {
		return "", err
	}
	return string(fromJson) + string(toJson) + strconv.FormatFloat(t.Fee, 'f', -1, 64), nil
}
