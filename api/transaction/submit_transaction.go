package transaction

import (
	"block_chain/block_structure/transaction"
	"block_chain/protocal/connection/content"
	"block_chain/utils"
	"encoding/json"
	"net"
)

func SubmitTransaction(Conn net.Conn, From []transaction.From, To []transaction.To, Fee float64, PublicKey string, Signature string) error {
	t := transaction.Transaction{From: From, To: To, Fee: Fee, PublicKey: PublicKey, Signature: Signature}
	strVal, err := t.ToString()
	if err != nil {
		return err
	}
	t.TransactionHash = utils.HashSHA256(strVal)
	request, err := json.Marshal(content.BroadcastTransaction{Transaction: t})
	if err != nil {
		return err
	}

	content.SendContent(Conn, string(request))
	return nil
}
