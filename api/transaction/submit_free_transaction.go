package transaction

import (
	"block_chain/block_structure/transaction"
	"block_chain/protocal/conseous"
	"errors"
	"net"
	"strconv"
)

func SubmitFreeTransaction(Conn net.Conn, Amount float64, Address string) error {
	if Amount > conseous.MasterAmount || Amount < 0 {
		return errors.New("amount can't large than " + strconv.Itoa(conseous.MasterAmount) + " or less than 0")
	}
	var f []transaction.From
	f = append(f, transaction.From{UTXOHash: conseous.MasterHash})

	var t []transaction.To
	t = append(t, transaction.To{Address: Address, Amount: Amount})

	return SubmitTransaction(Conn, f, t, 0, conseous.MasterPublicKey, conseous.MasterSignature)
}
