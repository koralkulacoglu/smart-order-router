package benchmarks

import (
	"container/heap"
	"sync"
)

type MutexOrderBook struct {
	sync.Mutex
	Bids *BidHeap
	Asks *AskHeap
}

func NewMutexOrderBook() *MutexOrderBook {
	bids := &BidHeap{}
	asks := &AskHeap{}
	heap.Init(bids)
	heap.Init(asks)
	return &MutexOrderBook{
		Bids: bids,
		Asks: asks,
	}
}

func (m *MutexOrderBook) AddOrders(bids, asks []OrderBookEntry) {
	m.Lock()
	defer m.Unlock()

	for _, b := range bids {
		heap.Push(m.Bids, b)
	}
	for _, a := range asks {
		heap.Push(m.Asks, a)
	}
}
