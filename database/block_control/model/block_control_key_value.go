package model

const (
	BlockControlKey        = "type"
	LastBlockKeyValue      = "LastBlock"
	CheckedBlockKeyValue   = "CheckedBlock"
	CandidateBlockKeyValue = "CandidateBlock"
)

type BlockControl struct {
	Type      string `bson:"type"`
	BlockHash string `bson:"block_hash"`
}
