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
