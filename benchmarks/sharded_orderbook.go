package benchmarks

import (
	"hash/fnv"
)

type ShardedOrderBook struct {
	shards []*MutexOrderBook
	mask   uint32
}

func NewShardedOrderBook(numShards int) *ShardedOrderBook {
	sob := &ShardedOrderBook{
		shards: make([]*MutexOrderBook, numShards),
		mask:   uint32(numShards - 1),
	}
	for i := 0; i < numShards; i++ {
		sob.shards[i] = NewMutexOrderBook()
	}
	return sob
}

func (sob *ShardedOrderBook) AddOrders(bids, asks []OrderBookEntry) {
	var symbol string
	if len(bids) > 0 {
		symbol = bids[0].Symbol
	} else {
		symbol = asks[0].Symbol
	}

	shardIndex := sob.getShardIndex(symbol)

	sob.shards[shardIndex].AddOrders(bids, asks)
}

func (sob *ShardedOrderBook) getShardIndex(symbol string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(symbol))
	return h.Sum32() & sob.mask
}
