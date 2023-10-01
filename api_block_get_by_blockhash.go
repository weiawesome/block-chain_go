package main

import (
	"block_chain/api"
	"fmt"
)

func main() {

	var BlockHash string

	BlockHash = "BlockHash"

	clientAPI := api.ClientAPI{ConnectNodeAddr: "127.0.0.1:8081"}
	err := clientAPI.Connect()
	if err != nil {
		return
	}
	defer func() {
		err := clientAPI.DisConnect()
		if err != nil {
			return
		}
	}()
	block, err := clientAPI.GetBlock(BlockHash)
	if err != nil {
		return
	}
	fmt.Println()
	fmt.Println("____________________________________________________________________________________________________")
	fmt.Println("BlockHash: ", block.Block.BlockHash)
	fmt.Println("BlockMeta: ", block.Block.BlockMeta.Content)
	fmt.Println("____________________________________________________________________________________________________")
	fmt.Println()
	fmt.Println("____________________________________________________________________________________________________")
	fmt.Println("BlockTop Information: ")
	fmt.Println("BlockTop-Version: ", block.Block.BlockTop.Version)
	fmt.Println("BlockTop-Height: ", block.Block.BlockTop.BlockHeight)
	fmt.Println("BlockTop-PreviousHash: ", block.Block.BlockTop.PreviousHash)
	fmt.Println("BlockTop-TimeStamp: ", block.Block.BlockTop.TimeStamp)
	fmt.Println("BlockTop-Nonce: ", block.Block.BlockTop.Nonce)
	fmt.Println("BlockTop-Difficulty: ", block.Block.BlockTop.Difficulty)
	fmt.Println("____________________________________________________________________________________________________")
	fmt.Println()
	fmt.Println("____________________________________________________________________________________________________")
	fmt.Println("BlockTransaction Information: ")
	for _, transaction := range block.Block.BlockTransactions {
		fmt.Println("TransactionHash: ", transaction.TransactionHash)
		fmt.Println("From: ")
		for _, from := range transaction.From {
			fmt.Println("	Address: ", from.UTXOHash)
			fmt.Println("	Index: ", from.Index)
		}
		fmt.Println("To: ")
		for _, to := range transaction.To {
			fmt.Println("	Address: ", to.Address)
			fmt.Println("	Amount: ", to.Amount)
		}
		fmt.Println("Fee: ", transaction.Fee)
		fmt.Println()
	}
	fmt.Println("____________________________________________________________________________________________________")
}
