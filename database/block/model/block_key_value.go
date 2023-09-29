package model

import "block_chain/block_structure/blockchain"

const (
	BlockKey       = "block_hash"
	BlockHeightKey = "block_height"
)

type BlockKeyValue struct {
	BlockHash   string           `bson:"block_hash"`
	BlockHeight int64            `bson:"block_height"`
	Block       blockchain.Block `bson:"block"`
}
