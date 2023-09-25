package utils

import (
	"context"
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

func InitClient(address string) error {
	clientOptions := options.Client().ApplyURI("mongodb://" + address)

	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	blockCollection := GetBlockCollection()
	_, err = blockCollection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    map[string]interface{}{BlockKey: 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}
	utxoCollection := GetUTXOCollection()
	_, err = utxoCollection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: map[string]interface{}{UTXOKey: 1, UTXOIndex: 1},
			Options: options.Index().
				SetUnique(true),
		},
	)
	if err != nil {
		return err
	}
	return err
}
func CloseMongoDb() error {
	return client.Disconnect(context.TODO())
}

func GetClient() *mongo.Client {
	return client
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
