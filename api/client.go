package api

import (
	BlockAPI "block_chain/api/block"
	"block_chain/api/content"
	TransactionAPI "block_chain/api/transaction"
	"block_chain/block_structure/transaction"
	"net"
)

type ClientAPI struct {
	ConnectNodeAddr string
	Conn            net.Conn
}

func (c *ClientAPI) GetBlock(BlockHash string) (content.ReturnBlock, error) {
	return BlockAPI.GetBlock(c.Conn, BlockHash)
}
func (c *ClientAPI) GetBlockByHeight(BlockHeight int64) ([]content.ReturnBlock, error) {
	return BlockAPI.GetBlockByHeight(c.Conn, BlockHeight)
}
func (c *ClientAPI) GetLastBlock() (content.ReturnBlock, error) {
	return BlockAPI.GetBlock(c.Conn, content.InitQuery)
}

func (c *ClientAPI) SubmitFreeTransaction(Amount float64, Address string) error {
	return TransactionAPI.SubmitFreeTransaction(c.Conn, Amount, Address)
}
func (c *ClientAPI) SubmitTransaction(From []transaction.From, To []transaction.To, Fee float64, PublicKey string, Signature string) error {
	return TransactionAPI.SubmitTransaction(c.Conn, From, To, Fee, PublicKey, Signature)
}

func (c *ClientAPI) Connect() error {
	conn, err := net.Dial("tcp", c.ConnectNodeAddr)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}
func (c *ClientAPI) DisConnect() error {
	err := c.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}
