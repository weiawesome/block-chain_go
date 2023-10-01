package block

import (
	"block_chain/api/content"
	"encoding/json"
	"net"
	"strings"
)

func GetBlock(Conn net.Conn, BlockHash string) (content.ReturnBlock, error) {
	request, err := json.Marshal(content.QueryBlock{QueryType: content.QueryByHash, BlockHash: BlockHash})
	if err != nil {
		return content.ReturnBlock{}, err
	}
	content.SendContent(Conn, string(request))
	totalResponse := ""
	buffer := make([]byte, content.BufferSize)

	for {
		_, err := Conn.Read(buffer)
		if err != nil {
			break
		}
		response := string(buffer)
		totalResponse += response
		if strings.ContainsAny(totalResponse, content.SuffixString) {
			totalResponse = totalResponse[:strings.Index(totalResponse, content.SuffixString)]
			break
		}
	}
	var block content.ReturnBlock
	err = json.Unmarshal([]byte(totalResponse), &block)
	return block, nil
}
