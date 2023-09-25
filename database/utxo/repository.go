package utxo

import (
	"block_chain/block_structure/blockchain"
	"block_chain/database/utxo/model"
	"block_chain/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func SetUTXO(transaction blockchain.BlockTransaction) error {
	UTXOCollection := utils.GetUTXOCollection()
	for i, to := range transaction.To {
		if _, err := GetUTXO(transaction.TransactionHash, i); err != nil {
			_, err = UTXOCollection.InsertOne(context.TODO(), model.UTXOKeyValue{TransactionHash: transaction.TransactionHash, Index: i, Amount: to.Amount, Address: to.Address, Spent: false})
			if err != nil {
				return err
			}
		}
	}
	for i, from := range transaction.From {
		if utxo, err := GetUTXO(from.UTXOHash, i); err != nil {
			_, err = UTXOCollection.InsertOne(context.TODO(), model.UTXOKeyValue{TransactionHash: transaction.TransactionHash, Index: i, Spent: true})
			if err != nil {
				return err
			}
		} else {
			utxo.Spent = true
			err := ReplaceUTXO(utxo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func ReplaceUTXO(value model.UTXOKeyValue) error {
	filter := bson.M{model.UTXOKey: value.TransactionHash, model.UTXOIndex: value.Index}
	UTXOCollection := utils.GetUTXOCollection()
	_, err := UTXOCollection.ReplaceOne(context.TODO(), filter, value)
	return err
}
func GetUTXO(TransactionHash string, Index int) (model.UTXOKeyValue, error) {
	UTXOCollection := utils.GetUTXOCollection()
	filter := bson.M{model.UTXOKey: TransactionHash, model.UTXOIndex: Index}

	var result model.UTXOKeyValue
	err := UTXOCollection.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}
