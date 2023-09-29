package connection

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"block_chain/utils"
	"encoding/json"
	"fmt"
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

func SendContent(conn net.Conn, val string) {
	fmt.Println(val)
	responseBytes := []byte(string(val) + SuffixString)

	_, err := conn.Write(responseBytes)
	if err != nil {
		utils.LogError(err.Error())
		return
	}
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func SentErrorMessage(conn net.Conn, message string) {
	response, err := json.Marshal(ErrorMessage{Error: message})
	if err != nil {
		utils.LogError(err.Error())
	}
	SendContent(conn, string(response))
}
