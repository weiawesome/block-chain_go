package connection

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"block_chain/database/block"
	"block_chain/database/block_control"
	"block_chain/initalize"
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

func InitBlocks(conn net.Conn) {
	flag := false
	request, err := json.Marshal(InitBlock{BlockHash: InitQuery})
	if err != nil {
		utils.LogError(err.Error())
		return
	}
	SendContent(conn, string(request))
	for {
		totalResponse := ""
		buffer := make([]byte, BufferSize)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				utils.LogError(err.Error())
				break
			}
			response := string(buffer[n:])
			totalResponse += response
			if strings.HasSuffix(totalResponse, SuffixString) {
				totalResponse = totalResponse[:len(totalResponse)-len(SuffixString)]
				break
			}
		}

		var data interface{}
		err := json.Unmarshal([]byte(totalResponse), &data)
		if err != nil {
			utils.LogError(err.Error())
			return
		}
		switch v := data.(type) {
		case ReturnBlock:
			if err := block.SetBlock(v.Block); err == nil {
				err := initalize.BuildUTXO(v.Block)
				if err != nil {
					return
				}
				if v.Block.BlockTop.PreviousHash == conseous.GenesisBlockPreviousHash {
					return
				}
				request, err := json.Marshal(InitBlock{BlockHash: v.Block.BlockTop.PreviousHash})
				if err != nil {
					utils.LogError(err.Error())
					return
				}
				SendContent(conn, string(request))
				if !flag {
					err := block_control.SetLastBlock(v.Block.BlockHash)
					if err != nil {
						return
					}
				}
				flag = true
			} else {
				return
			}
		default:
			SentErrorMessage(conn, "Unknown Request")
			return
		}
	}
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
			_, err := conn.Read(buffer)
			if err != nil {
				utils.LogError(err.Error())
				break
			}
			request := string(buffer)
			totalRequest += request
			if strings.ContainsAny(totalRequest, SuffixString) {
				totalRequest = totalRequest[:strings.Index(totalRequest, SuffixString)]
				fmt.Println("Success")
				break
			}
		}
		fmt.Println(totalRequest)
		var data interface{}
		err := json.Unmarshal([]byte(totalRequest), &data)
		if err != nil {
			utils.LogError(err.Error())
			return
		}
		switch v := data.(type) {
		case BroadcastTransaction:
			fmt.Println(v.Transaction.TransactionHash)
			transactionChannel <- v.Transaction
			break
		case BroadcastBlock:
			blockChannel <- v.Block
			break
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
			break
		case ErrorMessage:
			fmt.Println(v.Error)
			utils.LogError(v.Error)
			break
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
			fmt.Println(val)
			connections = append(connections, val)
		case b := <-BroadcastBlockChannel:
			response, err := json.Marshal(BroadcastBlock{Block: b})
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
			fmt.Println(t.TransactionHash)
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
		//go ReceiveReplyClient(conn, TransactionChannel, BlockChannel)
		//go InitBlocks(conn)
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
		fmt.Println(conn.RemoteAddr())
		go ReceiveReplyClient(conn, TransactionChannel, BlockChannel)
		ConnectionChannel <- conn
	}
}
