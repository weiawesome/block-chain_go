package initalize

import (
	"block_chain/database/block"
	"block_chain/database/block_control"
	"block_chain/protocal/connection"
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"encoding/json"
	"net"
	"strings"
)

func InitBlocks(conn net.Conn) {
	flag := false
	request, err := json.Marshal(connection.InitBlock{BlockHash: connection.InitQuery})
	if err != nil {
		utils.LogError(err.Error())
		return
	}
	connection.SendContent(conn, string(request))
	for {
		totalResponse := ""
		buffer := make([]byte, connection.BufferSize)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				utils.LogError(err.Error())
				break
			}
			response := string(buffer[n:])
			totalResponse += response
			if strings.HasSuffix(totalResponse, connection.SuffixString) {
				totalResponse = totalResponse[:len(totalResponse)-len(connection.SuffixString)]
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
		case connection.ReturnBlock:
			if err := block.SetBlock(v.Block); err == nil {
				err := BuildUTXO(v.Block)
				if err != nil {
					return
				}
				if v.Block.BlockTop.PreviousHash == conseous.GenesisBlockPreviousHash {
					return
				}
				request, err := json.Marshal(connection.InitBlock{BlockHash: v.Block.BlockTop.PreviousHash})
				if err != nil {
					utils.LogError(err.Error())
					return
				}
				connection.SendContent(conn, string(request))
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
			response, err := json.Marshal(connection.ErrorMessage{Error: "Unknown request"})
			if err != nil {
				utils.LogError(err.Error())
				continue
			}
			connection.SendContent(conn, string(response))
			return
		}
	}
}
