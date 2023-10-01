package connection

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"block_chain/database/block"
	"block_chain/database/block_control"
	"block_chain/initalize"
	"block_chain/protocal/connection/content"
	"block_chain/protocal/conseous"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

func InitBlocks(conn net.Conn) {
	flag := false
	request, err := json.Marshal(content.QueryBlock{QueryType: content.QueryByHash, BlockHash: content.InitQuery})
	if err != nil {
		fmt.Println("Error found ", err.Error(), " in BN")
		return
	}
	content.SendContent(conn, string(request))
	for {
		totalResponse := ""
		buffer := make([]byte, content.BufferSize)
		connFlag := false
		for {
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error found ", err.Error(), " in BN")
				connFlag = true
				break
			}
			response := string(buffer)
			totalResponse += response
			if strings.ContainsAny(totalResponse, content.SuffixString) {
				totalResponse = totalResponse[:strings.Index(totalResponse, content.SuffixString)]
				break
			}
		}
		if connFlag {
			break
		}

		var data content.ReturnBlock
		err := json.Unmarshal([]byte(totalResponse), &data)
		if err == nil && data.Block.BlockHash != "" {
			if err := block.SetBlock(data.Block); err == nil {
				err := initalize.BuildUTXO(data.Block)
				if err != nil {
					return
				}
				if data.Block.BlockTop.PreviousHash == conseous.GenesisBlockPreviousHash {
					return
				}
				request, err := json.Marshal(content.QueryBlock{QueryType: content.QueryByHash, BlockHash: data.Block.BlockTop.PreviousHash})
				if err != nil {
					fmt.Println("Error found ", err.Error(), " in BN")
					return
				}
				content.SendContent(conn, string(request))
				if !flag {
					err := block_control.SetLastBlock(data.Block.BlockHash)
					if err != nil {
						return
					}
				}
				flag = true
			} else {
				return
			}
			continue
		} else {
			content.SentErrorMessage(conn, "Unknown Request")
			return
		}
	}
}

func ReceiveReplyClient(conn net.Conn, transactionChannel chan transaction.Transaction, blockChannel chan blockchain.Block) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	buffer := make([]byte, content.BufferSize)

	for {
		flag := false
		totalRequest := ""
		for {
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error found ", err.Error(), " in BN")
				flag = true
				break
			}
			request := string(buffer)
			totalRequest += request
			if strings.ContainsAny(totalRequest, content.SuffixString) {
				totalRequest = totalRequest[:strings.Index(totalRequest, content.SuffixString)]
				break
			}
		}
		if flag {
			break
		}
		fmt.Println("Get message: ", totalRequest, " in BN")

		var broadcastTransactionData content.BroadcastTransaction
		var broadcastBlockData content.BroadcastBlock
		var queryBlockData content.QueryBlock
		var queryBlockByHeightData content.QueryBlockByHeight
		var errorMessageData content.ErrorMessage
		var err error

		err = json.Unmarshal([]byte(totalRequest), &broadcastTransactionData)
		if err == nil && broadcastTransactionData.Transaction.TransactionHash != "" {
			fmt.Println("Get new Transaction: ", broadcastTransactionData.Transaction.TransactionHash, " in BN")
			fmt.Println("Send new Transaction: ", broadcastTransactionData.Transaction.TransactionHash, " in BN")
			fmt.Println()
			transactionChannel <- broadcastTransactionData.Transaction
			continue
		}
		err = json.Unmarshal([]byte(totalRequest), &broadcastBlockData)
		if err == nil && broadcastBlockData.Block.BlockHash != "" {
			fmt.Println("Get new Block: ", broadcastBlockData.Block.BlockHash, " in BN")
			fmt.Println("Send new Block: ", broadcastBlockData.Block.BlockHash, " in BN")
			fmt.Println()
			blockChannel <- broadcastBlockData.Block
			continue
		}
		err = json.Unmarshal([]byte(totalRequest), &queryBlockData)
		if err == nil && queryBlockData.BlockHash != "" {
			fmt.Println("Get new Query Block: ", queryBlockData.BlockHash, " in BN")
			if queryBlockData.BlockHash == content.InitQuery {
				lastBlock, err := block_control.GetLastBlock()
				if err != nil {
					content.SentErrorMessage(conn, "Error to get last block")
					fmt.Println("Error found ", err.Error(), " in BN")
					continue
				}
				Block, err := block.GetBlock(lastBlock)
				fmt.Println("Sent Last Block ", Block.BlockHash, " in BN")
				fmt.Println()
				if err != nil {
					content.SentErrorMessage(conn, "Error to get block")
					fmt.Println("Error found ", err.Error(), " in BN")
					continue
				}
				response, err := json.Marshal(content.ReturnBlock{Block: Block})
				if err != nil {
					content.SentErrorMessage(conn, "Error to send block")
					fmt.Println("Error found ", err.Error(), " in BN")
					continue
				}
				content.SendContent(conn, string(response))
			} else {
				Block, err := block.GetBlock(queryBlockData.BlockHash)
				if err != nil {
					content.SentErrorMessage(conn, "Error to get block")
					fmt.Println("Error found ", err.Error(), " in BN")
					continue
				}
				response, err := json.Marshal(content.ReturnBlock{Block: Block})
				if err != nil {
					content.SentErrorMessage(conn, "Error to send block")
					fmt.Println("Error found ", err.Error(), " in BN")
					continue
				}
				fmt.Println("Sent Block ", Block.BlockHash, " in BN")
				fmt.Println()
				content.SendContent(conn, string(response))
			}
			continue
		}
		err = json.Unmarshal([]byte(totalRequest), &queryBlockByHeightData)
		if err == nil && queryBlockByHeightData.QueryType != "" {
			fmt.Println("Get new Query Block by height: ", queryBlockByHeightData.BlockHeight, " in BN")
			blocks, err := block.GetBlockByBlockHeight(queryBlockByHeightData.BlockHeight)
			if err != nil {
				content.SentErrorMessage(conn, err.Error())
				continue
			}
			var returnBlocks []content.ReturnBlock
			for _, b := range blocks {
				returnBlocks = append(returnBlocks, content.ReturnBlock{Block: b})
			}
			response, err := json.Marshal(returnBlocks)
			if err != nil {
				content.SentErrorMessage(conn, "Error to send block")
				fmt.Println("Error found ", err.Error(), " in BN")
				continue
			}
			fmt.Println("Sent Block by height in BN")
			fmt.Println()
			content.SendContent(conn, string(response))
			continue
		}
		err = json.Unmarshal([]byte(totalRequest), &errorMessageData)
		if err == nil && errorMessageData.Error != "" {
			fmt.Println("Get Error Message ", errorMessageData.Error, " in BN")
			fmt.Println("Sent Error Message ", "Unknown request", " in BN")
			fmt.Println()
			continue
		}
		content.SentErrorMessage(conn, "Unknown request")
	}
}
func CommunicateClient(ConnectionChannel chan net.Conn, BroadcastTransactionChannel chan transaction.Transaction, BroadcastBlockChannel chan blockchain.Block) {
	var connections []net.Conn
	for {
		select {
		case val := <-ConnectionChannel:
			fmt.Println("Get new connection ", val.RemoteAddr(), " in BN")
			connections = append(connections, val)
		case b := <-BroadcastBlockChannel:
			fmt.Println("Get new BroadcastBlock ", b.BlockHash, " in BN")
			fmt.Println("Sent new BroadcastBlock ", b.BlockHash, " in BN")
			fmt.Println()
			response, err := json.Marshal(content.BroadcastBlock{Block: b})
			for _, connection := range connections {
				if err != nil {
					content.SentErrorMessage(connection, "Error to send block")
					fmt.Println("Error found ", err.Error(), " in BN")
					continue
				}
				content.SendContent(connection, string(response))
			}
		case t := <-BroadcastTransactionChannel:
			fmt.Println("Get new Transaction ", t.TransactionHash, " in BN")
			fmt.Println("Sent new Transaction ", t.TransactionHash, " in BN")
			fmt.Println()
			response, err := json.Marshal(content.BroadcastTransaction{Transaction: t})
			for _, connection := range connections {
				if err != nil {
					content.SentErrorMessage(connection, "Error to send transaction")
					fmt.Println("Error found ", err.Error(), " in BN")
					continue
				}
				content.SendContent(connection, string(response))
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
		go ReceiveReplyClient(conn, TransactionChannel, BlockChannel)
		go InitBlocks(conn)
		ConnectionChannel <- conn
	}

	listener, err := net.Listen("tcp", NodeAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(listener)

	fmt.Println("Started to build node in BN")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error found ", err.Error(), " in BN")
			continue
		}
		fmt.Println("Get new node ", conn.RemoteAddr(), " in BN")
		go ReceiveReplyClient(conn, TransactionChannel, BlockChannel)
		ConnectionChannel <- conn
	}
}
