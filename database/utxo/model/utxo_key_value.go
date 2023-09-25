package model

const (
	UTXOKey   = "transaction_hash"
	UTXOIndex = "index"
)

type UTXOKeyValue struct {
	TransactionHash string  `bson:"transaction_hash"`
	Index           int     `bson:"index"`
	Amount          float64 `bson:"amount"`
	Spent           bool    `bson:"spent"`
	Address         string  `bson:"address"`
}
