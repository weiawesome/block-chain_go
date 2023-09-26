package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DataBase        = "DataBase"
	BlockCollection = "Block"
	UTXOCollection  = "UTXO"
	BLockControl    = "BlockControl"
)
const (
	BlockKey = "block_hash"
)
const (
	UTXOKey   = "transaction_hash"
	UTXOIndex = "index"
)

var client *mongo.Client

func BlockCollectInit() error {
	blockCollection := GetBlockCollection()
	_, err := blockCollection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    map[string]interface{}{BlockKey: 1},
			Options: options.Index().SetUnique(true),
		},
	)
	return err
}

func UTXOCollectInit() error {
	utxoCollection := GetUTXOCollection()
	_, err := utxoCollection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: UTXOKey, Value: 1},
				{Key: UTXOIndex, Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	)
	return err
}

func InitClient(address string) error {
	var err error

	clientOptions := options.Client().ApplyURI("mongodb://" + address)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	err = BlockCollectInit()
	if err != nil {
		return err
	}

	err = UTXOCollectInit()
	if err != nil {
		return err
	}

	return nil
}

func CloseMongoDb() error {
	return client.Disconnect(context.TODO())
}

func GetBlockCollection() *mongo.Collection {
	return client.Database(DataBase).Collection(BlockCollection)
}
func GetUTXOCollection() *mongo.Collection {
	return client.Database(DataBase).Collection(UTXOCollection)
}
func GetBlockControlCollection() *mongo.Collection {
	return client.Database(DataBase).Collection(BLockControl)
}
