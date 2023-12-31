# Block Chain 區塊鏈
> This is a project that I try to build a blockchain by myself.
> 
> Using go lang to build a blockchain.

## How to download
```shell
git clone https://github.com/weiawesome/block-chain_go.git
```
## How to start (Miners)
(Ensure the computer have been installed go)
(Ensure the mongodb is working)
```shell
go run main.go
```
![](https://hackmd.io/_uploads/HJ72MCvl6.png)
* Success to build a node in 127.0.0.1:8080

### How to edit the file
[Edit the main.go](main.go)
```go
//This mean the miner's node address
NodeAddr := "127.0.0.1:8080"

//This mean can connect to other node
NodeAddresses := []string{}

//This mean the mongodb's address
DbAddress:="localhost:27017"

//This mean the miners(goroutines) to mine. 
Miners := 1
```
### Easy to start a mongodb by docker
Make a mongodb container and port-forward to localhost
```shell
docker run --name mongodb -p 27017:27017 -d mongo
```
### How to connect with exist node
[Edit the main.go](main.go)
```go
//This mean the miner's node address(choose another port)
NodeAddr := "127.0.0.1:8081"

//This mean can connect to other node(choose exist node)
NodeAddresses := []string{"127.0.0.1:8080"}

//This mean the mongodb's address(choose db in other port)
DbAddress:="localhost:27018"
```
![](https://hackmd.io/_uploads/Hywt5l_xp.png)
* Success to connect with 127.0.0.1:8080
* Success to build and update the database
* Success to build a node in 127.0.0.1:8081 and connect with 127.0.0.1:8080 

## How to get block
(Ensure the Miners is working)
### Get block by block_hash
```shell
go run api_block_get_by_block_hash.go
```
![](https://hackmd.io/_uploads/rywjAJ_xp.png)
![](https://hackmd.io/_uploads/B1Ixkldep.png)
* Success to get the block with specific block hash

### Edit Get block by block_hash
[Edit the api_block_get_by_block_hash.go](api_block_get_by_block_hash.go)
```go
//This is the query block_hash
var BlockHash string
BlockHash = "BlockHash"

//This is the connect node address
ConnectAddr:="127.0.0.1:8080"
```
### Get block by block_height
```shell
go run api_block_get_by_block_height.go
```

![](https://hackmd.io/_uploads/BJF1bl_eT.png)
![](https://hackmd.io/_uploads/Sk-bbl_ea.png)
* Success to get the block which block height equal to 0(Genesis block)

### Edit Get block by block_height
[Edit the api_block_get_by_block_height.go](api_block_get_by_block_height.go)
```go
//This is the query block_height
var BlockHeight int64
BlockHeight = 0

//This is the connect node address
ConnectAddr := "127.0.0.1:8080"
```

### Get the last block 
```shell
go run api_block_get_last.go
```
![](https://hackmd.io/_uploads/rywjAJ_xp.png)
![](https://hackmd.io/_uploads/B1Ixkldep.png)
* Success to get latest block

### Edit Get the last block
[Edit the api_block_get_last.go](api_block_get_last.go)
```go
//This is the connect node address
ConnectAddr := "127.0.0.1:8080"
```

## How to submit transaction
(Ensure the Miners is working)
### Submit the transaction
```shell
go run api_transaction_submit.go
```
![](https://hackmd.io/_uploads/r1TRjyulp.png)
* Success to submit a transaction into blockchain
* Success to show the transaction's detail information

### Edit Submit the transaction
[Edit the api_transaction_submit.go](api_transaction_submit.go)
```go
//This is UTXOHash to prove the asset
UTXOHash = "UTXOHash"

//This is to use with UTXOHash to prove the asset
Index = 0

//This is the address the transaction transfer
Address = "Address"

//This is the amount the transaction transfer
Amount = 0

//This is the fee to miner in the transaction
Fee = 0

//This is public key for sender
PublicKey = "PublicKey"

//This is private key for sender
PrivateKey = "PrivateKey"

//This is the connect node address
ConnectAddr := "127.0.0.1:8080"
```

### Submit the Free transaction
```shell
# submit a transaction with specific address
go run api_transaction_submit_free.go
# submit a transaction with generating an address
go run api_transaction_submit_free_random_addr.go 
```
![](https://hackmd.io/_uploads/r1qxEyueT.png)
* Success to submit a transaction into blockchain
* Success to gain the free asset
 
### Edit Submit the Free transaction
[Edit the api_transaction_submit_free.go](api_transaction_submit_free.go)
```go
//This is the amount give to the address 
var Amount float64
Amount=0

//This is the address to be given
var Addr string
Addr = "Address"

//This is the connect node address
ConnectAddr := "127.0.0.1:8080"
```
[Edit the api_transaction_submit_free_random_addr.go](api_transaction_submit_free_random_addr.go)
```go
//This is the amount give to the address 
var Amount float64
Amount=0

//This is the connect node address
ConnectAddr := "127.0.0.1:8080"
```