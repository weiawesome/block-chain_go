package model

import "block_chain/block_structure/blockchain"

const (
	BlockKey = "block_hash"
)

type BlockKeyValue struct {
	BlockHash string           `bson:"block_hash"`
	Block     blockchain.Block `bson:"block"`
}

type BlockControl struct {
	Type      string `bson:"type"`
	BlockHash string `bson:"block_hash"`
}
