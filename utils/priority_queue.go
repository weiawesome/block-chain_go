package utils

import (
	"block_chain/block_structure/blockchain"
	"block_chain/block_structure/transaction"
	"container/heap"
	"fmt"
)

type PriorityQueue []*blockchain.BlockTransaction

func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i].Fee < (*pq)[j].Fee
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*blockchain.BlockTransaction)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func test() {
	pq := make(PriorityQueue, 0)
	item1 := &transaction.Transaction{}
	heap.Push(&pq, item1)
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*transaction.Transaction)
		fmt.Print(item.Fee)
	}
}
