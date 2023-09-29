package refresh_db

import (
	"block_chain/block_structure/blockchain"
	"block_chain/database/block"
	"block_chain/database/block_control"
	"block_chain/database/utxo"
	"block_chain/protocal/conseous"
	"fmt"
)

func RefreshDb(RefreshBlockChannel chan blockchain.Block, CompleteBlockChannel chan blockchain.Block) {
	for {
		select {
		case rb := <-RefreshBlockChannel:
			fmt.Println("Receive Block", rb.BlockHash, " in RDB")
			fmt.Println("Genesis Block ", rb.BlockTop.PreviousHash == conseous.GenesisBlockPreviousHash, " in RDB")
			if rb.BlockTop.PreviousHash == conseous.GenesisBlockPreviousHash {
				fmt.Println("Block is GenesisBlock", rb.BlockHash, " in RDB")
				lastHash, err := block_control.GetLastBlock()
				if err != nil {
					err := block_control.SetLastBlock(rb.BlockHash)
					if err != nil {
						continue
					}
					fmt.Println("Last block replace with", rb.BlockHash, " in RDB")
				} else if lastHash == conseous.GenesisBlockPreviousHash {
					err := block_control.SetLastBlock(rb.BlockHash)
					if err != nil {
						continue
					}
					fmt.Println("Last block replace with", rb.BlockHash, " in RDB")
				}
				err = block_control.SetCandidateBlock(rb.BlockHash)
				if err != nil {
					continue
				}
				err = block.SetBlock(rb)
				if err != nil {
					continue
				}
				for _, transaction := range rb.BlockTransactions {
					err := utxo.SetUTXO(transaction)
					if err != nil {
						continue
					}
				}
				fmt.Println("Sent Block in RDB")
				fmt.Println()
				CompleteBlockChannel <- rb
				continue
			} else {
				fmt.Println("Block is not Genesis in RDB")
				_, err := block.GetBlock(rb.BlockTop.PreviousHash)
				if err != nil {
					fmt.Println("Previous Block not exist", rb.BlockTop.PreviousHash, " in RDB")
					continue
				}
			}
			lastBlockHash, err := block_control.GetLastBlock()
			if err != nil {
				continue
			}
			fmt.Println("LastBlock", lastBlockHash, " in RDB")
			lastBlock, err := block.GetBlock(lastBlockHash)
			if err != nil {
				continue
			}
			fmt.Println("Success get LastBlock in RDB")
			if rb.BlockTop.BlockHeight > lastBlock.BlockTop.BlockHeight-conseous.BlockChecked {
				fmt.Println("In checked Block", rb.BlockHash, " in RDB")
				err := block_control.SetCandidateBlock(rb.BlockHash)
				if err != nil {
					continue
				}
				for _, transaction := range rb.BlockTransactions {
					err := utxo.SetUTXO(transaction)
					if err != nil {
						continue
					}
				}
				if rb.BlockTop.BlockHeight > lastBlock.BlockTop.BlockHeight {
					fmt.Println("Higher Block", rb.BlockHash, " in RDB")
					err := block_control.SetLastBlock(rb.BlockHash)
					if err != nil {
						continue
					}
					lastBlock = rb
					if lastBlock.BlockTop.BlockHeight-conseous.BlockChecked >= conseous.GenesisBlockHeight {
						tmpBlock := lastBlock
						for i := 0; i < conseous.BlockChecked; i++ {
							tmpBlock, err = block.GetBlock(tmpBlock.BlockTop.PreviousHash)
							if err != nil {
								continue
							}
						}
						err := block_control.DeleteCandidateBlock(tmpBlock.BlockHash)
						if err != nil {
							continue
						}
						candidateBlocks, err := block_control.GetCandidateBlock()
						if err != nil {
							continue
						}
						for _, candidateBlock := range candidateBlocks {
							b, err := block.GetBlock(candidateBlock)
							if err != nil {
								return
							}
							if b.BlockTop.BlockHeight == tmpBlock.BlockTop.BlockHeight+1 {
								if b.BlockTop.PreviousHash != tmpBlock.BlockHash {
									err := block_control.DeleteCandidateBlock(b.BlockHash)
									if err != nil {
										continue
									}
									for _, transaction := range b.BlockTransactions {
										err := utxo.ReverseUTXO(transaction)
										if err != nil {
											continue
										}
									}
									err = block.DeleteBlock(b.BlockHash)
									if err != nil {
										continue
									}
								}
							} else if b.BlockTop.BlockHeight == tmpBlock.BlockTop.BlockHeight+1 {
								err := block_control.DeleteCandidateBlock(b.BlockHash)
								if err != nil {
									continue
								}
								for _, transaction := range b.BlockTransactions {
									err := utxo.ReverseUTXO(transaction)
									if err != nil {
										continue
									}
								}
								err = block.DeleteBlock(b.BlockHash)
								if err != nil {
									continue
								}
							}

						}
						err = block_control.SetCheckedBlock(tmpBlock.BlockHash)
						if err != nil {
							continue
						}
					}
				}
				fmt.Println("Sent Complete Block not Genesis in RDB")
				fmt.Println()
				CompleteBlockChannel <- rb
			}

		default:
			continue
		}
	}
}
