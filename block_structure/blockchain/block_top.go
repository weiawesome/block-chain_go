package blockchain

import "strconv"

type BlockTop struct {
	Version      int32  `json:"version"`
	PreviousHash string `json:"previousHash"`
	TimeStamp    int64  `json:"timeStamp"`
	Difficulty   int64  `json:"difficulty"`
	Nonce        uint32 `json:"nonce"`
	MarkleRoot   string `json:"markleRoot"`
}

func (bt BlockTop) ToString() string {
	return string(bt.Version) +
		bt.PreviousHash +
		strconv.FormatInt(bt.TimeStamp, 10) +
		strconv.FormatInt(bt.Difficulty, 10) +
		strconv.FormatUint(uint64(bt.Nonce), 10) +
		bt.MarkleRoot
}
