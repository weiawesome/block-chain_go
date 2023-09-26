package block

import (
	"block_chain/block_structure/blockchain"
	"block_chain/database/block/model"
	"block_chain/protocal/conseous"
	"block_chain/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func GetBlock(Hash string) (blockchain.Block, error) {
	BlockCollection := utils.GetBlockCollection()
	filter := bson.M{model.BlockKey: Hash}

	var result model.BlockKeyValue
	err := BlockCollection.FindOne(context.TODO(), filter).Decode(&result)

	return result.Block, err
}

func SetBlock(Block blockchain.Block) error {
	BlockCollection := utils.GetBlockCollection()
	_, err := BlockCollection.InsertOne(context.TODO(), model.BlockKeyValue{BlockHash: Block.BlockHash, Block: Block})
	return err
}

func DeleteBlock(Hash string) error {
	BlockCollection := utils.GetBlockCollection()

	filter := bson.M{model.BlockKey: Hash}

	_, err := BlockCollection.DeleteOne(context.TODO(), filter)

	return err
}

func TransactionInCheckedBlocks(BlockHash string, TransactionHash string) bool {
	tmpHash := BlockHash
	for i := 0; i < conseous.BlockChecked; i++ {
		block, err := GetBlock(tmpHash)
		if err != nil {
			return true
		}
		for _, transaction := range block.BlockTransactions {
			if transaction.TransactionHash == TransactionHash {
				return false
			}
		}
	}
	return true
}
