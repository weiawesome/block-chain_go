package model

const (
	BlockControlKey   = "type"
	LastBlockKeyValue = "LastBlock"
)

type BlockControl struct {
	Type      string `bson:"type"`
	BlockHash string `bson:"block_hash"`
}
