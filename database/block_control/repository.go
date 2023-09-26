package block_control

import (
	"block_chain/database/block_control/model"
	"block_chain/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetLastBlock(Hash string) error {
	BlockControlCollection := utils.GetBlockControlCollection()

	filter := bson.M{model.BlockControlKey: model.LastBlockKeyValue}
	option := options.Replace().SetUpsert(true)

	_, err := BlockControlCollection.ReplaceOne(context.TODO(), filter, model.BlockControl{Type: model.LastBlockKeyValue, BlockHash: Hash}, option)
	return err
}
func GetLastBlock() (string, error) {
	BlockControlCollection := utils.GetBlockControlCollection()

	filter := bson.M{model.BlockControlKey: model.LastBlockKeyValue}
	var lastBlock model.BlockControl
	err := BlockControlCollection.FindOne(context.TODO(), filter).Decode(lastBlock)
	return lastBlock.BlockHash, err
}

func SetCheckedBlock(Hash string) error {
	BlockControlCollection := utils.GetBlockControlCollection()

	filter := bson.M{model.BlockControlKey: model.CheckedBlockKeyValue}
	option := options.Replace().SetUpsert(true)

	_, err := BlockControlCollection.ReplaceOne(context.TODO(), filter, model.BlockControl{Type: model.LastBlockKeyValue, BlockHash: Hash}, option)
	return err
}
func GetCheckedBlock() (string, error) {
	BlockControlCollection := utils.GetBlockControlCollection()

	filter := bson.M{model.BlockControlKey: model.CheckedBlockKeyValue}
	var lastBlock model.BlockControl
	err := BlockControlCollection.FindOne(context.TODO(), filter).Decode(lastBlock)
	return lastBlock.BlockHash, err
}

func SetCandidateBlock(Hash string) error {
	BlockControlCollection := utils.GetBlockControlCollection()

	_, err := BlockControlCollection.InsertOne(context.TODO(), model.BlockControl{Type: model.LastBlockKeyValue, BlockHash: Hash})

	return err
}
func GetCandidateBlock() ([]string, error) {
	var results []string

	BlockControlCollection := utils.GetBlockControlCollection()

	filter := bson.M{model.BlockControlKey: model.CandidateBlockKeyValue}

	cursor, err := BlockControlCollection.Find(context.TODO(), filter)
	if err != nil {
		return results, err
	}

	for cursor.Next(context.TODO()) {
		var result model.BlockControl
		if err := cursor.Decode(&result); err != nil {
			return results, err
		}
		results = append(results, result.BlockHash)
	}
	return results, err
}
func DeleteCandidateBlock(Hash string) error {
	BlockControlCollection := utils.GetBlockControlCollection()

	filter := bson.M{model.CandidateBlockKeyValue: Hash}

	_, err := BlockControlCollection.DeleteOne(context.TODO(), filter)

	return err
}
