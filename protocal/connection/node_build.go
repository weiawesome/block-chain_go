package connection

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"block_chain/database/block"
	"block_chain/database/block_control"
	"block_chain/utils"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

func SentErrorMessage(conn net.Conn, message string) {
	response, err := json.Marshal(ErrorMessage{Error: message})
	if err != nil {
		utils.LogError(err.Error())
	}
	SendContent(conn, string(response))
}
func ReceiveReplyClient(conn net.Conn, transactionChannel chan transaction.Transaction, blockChannel chan blockchain.Block) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	buffer := make([]byte, BufferSize)

	for {
		totalRequest := ""
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				utils.LogError(err.Error())
				break
			}
			request := string(buffer[n:])
			totalRequest += request
			if strings.HasSuffix(totalRequest, SuffixString) {
				totalRequest = totalRequest[:len(totalRequest)-len(SuffixString)]
				break
			}
		}

		var data interface{}
		err := json.Unmarshal([]byte(totalRequest), &data)
		if err != nil {
			utils.LogError(err.Error())
			return
		}
		switch v := data.(type) {
		case BroadcastTransaction:
			transactionChannel <- v.Transaction
		case BroadcastBlock:
			blockChannel <- v.Block
		case InitBlock:
			if v.BlockHash == InitQuery {
				lastBlock, err := block_control.GetLastBlock()
				if err != nil {
					SentErrorMessage(conn, "Error to get last block")
					utils.LogError(err.Error())
					continue
				}
				Block, err := block.GetBlock(lastBlock)
				if err != nil {
					SentErrorMessage(conn, "Error to get block")
					utils.LogError(err.Error())
					continue
				}
				response, err := json.Marshal(ReturnBlock{Block: Block})
				if err != nil {
					SentErrorMessage(conn, "Error to send block")
					utils.LogError(err.Error())
					continue
				}
				SendContent(conn, string(response))
			} else {
				Block, err := block.GetBlock(v.BlockHash)
				if err != nil {
					SentErrorMessage(conn, "Error to get block")
					utils.LogError(err.Error())
					continue
				}
				response, err := json.Marshal(ReturnBlock{Block: Block})
				if err != nil {
					SentErrorMessage(conn, "Error to send block")
					utils.LogError(err.Error())
					continue
				}
				SendContent(conn, string(response))
			}
		case ErrorMessage:
			utils.LogError(v.Error)
		default:
			SentErrorMessage(conn, "Unknown request")
		}
	}
}
func CommunicateClient(ConnectionChannel chan net.Conn, BroadcastTransactionChannel chan transaction.Transaction, BroadcastBlockChannel chan blockchain.Block) {
	var connections []net.Conn
	for {
		select {
		case val := <-ConnectionChannel:
			connections = append(connections, val)
		case b := <-BroadcastBlockChannel:
			response, err := json.Marshal(BroadcastBlock{Block: b})
			fmt.Println(b.BlockHash)
			for _, connection := range connections {
				if err != nil {
					SentErrorMessage(connection, "Error to send block")
					utils.LogError(err.Error())
					continue
				}
				SendContent(connection, string(response))
			}
		case t := <-BroadcastTransactionChannel:
			response, err := json.Marshal(BroadcastTransaction{Transaction: t})
			for _, connection := range connections {
				if err != nil {
					SentErrorMessage(connection, "Error to send transaction")
					utils.LogError(err.Error())
					continue
				}
				SendContent(connection, string(response))
			}
		default:
			continue
		}

	}
}

func BuildNode(NodeAddr string, NodeAddresses []string, TransactionChannel chan transaction.Transaction, BlockChannel chan blockchain.Block, BroadcastTransactionChannel chan transaction.Transaction, BroadcastBlockChannel chan blockchain.Block) {
	ConnectionChannel := make(chan net.Conn)
	go CommunicateClient(ConnectionChannel, BroadcastTransactionChannel, BroadcastBlockChannel)

	for _, addr := range NodeAddresses {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			continue
		}
		//go initalize.InitBlocks(conn)
		ConnectionChannel <- conn
	}

	listener, err := net.Listen("tcp", NodeAddr)
	if err != nil {
		panic(err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.LogError(err.Error())
			continue
		}
		go ReceiveReplyClient(conn, TransactionChannel, BlockChannel)
		ConnectionChannel <- conn
	}
}
