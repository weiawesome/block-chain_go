package connection

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"block_chain/utils"
	"net"
)

type BroadcastTransaction struct {
	Transaction transaction.Transaction
}

type BroadcastBlock struct {
	Block blockchain.Block `json:"block"`
}

const InitQuery = "InitQuery"

type InitBlock struct {
	BlockHash string `json:"blockHash"`
}
type ReturnBlock struct {
	Block blockchain.Block `json:"block"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func SendContent(conn net.Conn, val string) {
	responseBytes := []byte(string(val) + SuffixString)

	_, err := conn.Write(responseBytes)
	if err != nil {
		utils.LogError(err.Error())
		return
	}
}
