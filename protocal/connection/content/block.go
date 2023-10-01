package content

import "block_chain/block_structure/blockchain"

type BroadcastBlock struct {
	Block blockchain.Block `json:"block"`
}
