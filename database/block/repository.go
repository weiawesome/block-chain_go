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

	var finalResult model.BlockKeyValue
	cursor, err := BlockCollection.Find(context.TODO(), filter)
	for cursor.Next(context.TODO()) {
		var result model.BlockKeyValue
		if err := cursor.Decode(&result); err != nil {
			return finalResult.Block, err
		}
		finalResult = result
	}
	return finalResult.Block, err
}

func GetBlockByBlockHeight(Height int64) ([]blockchain.Block, error) {
	var finalResult []blockchain.Block
	BlockCollection := utils.GetBlockCollection()
	filter := bson.M{model.BlockHeightKey: Height}
	cursor, err := BlockCollection.Find(context.TODO(), filter)
	for cursor.Next(context.TODO()) {
		var result model.BlockKeyValue
		if err := cursor.Decode(&result); err != nil {
			return finalResult, err
		}
		finalResult = append(finalResult, result.Block)
	}

	return finalResult, err
}

func SetBlock(Block blockchain.Block) error {
	BlockCollection := utils.GetBlockCollection()
	_, err := BlockCollection.InsertOne(context.TODO(), model.BlockKeyValue{BlockHash: Block.BlockHash, BlockHeight: Block.BlockTop.BlockHeight, Block: Block})
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

func CheckGenerateSpeed(Height int64) (bool, error) {
	blockLower, errLower := GetBlockByBlockHeight(Height - 1)
	if errLower != nil {
		return false, errLower
	}
	blockUpper, errUpper := GetBlockByBlockHeight(Height - conseous.DifficultyCycle)
	if errUpper != nil {
		return false, errUpper
	}
	return blockLower[0].BlockTop.TimeStamp-blockUpper[0].BlockTop.TimeStamp > conseous.AverageBlockGenerateTime*conseous.DifficultyCycle, nil
}
