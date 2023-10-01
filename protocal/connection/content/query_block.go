package content

import "block_chain/block_structure/blockchain"

const (
	InitQuery     = "InitQuery"
	QueryByHash   = "QueryByHash"
	QueryByHeight = "QueryByHeight"
)

type QueryBlock struct {
	QueryType string `json:"type"`
	BlockHash string `json:"blockHash"`
}

type QueryBlockByHeight struct {
	QueryType   string `json:"type"`
	BlockHeight int64  `json:"block_height"`
}

type ReturnBlock struct {
	Block blockchain.Block `json:"block"`
}
