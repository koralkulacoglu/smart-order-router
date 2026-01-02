package benchmarks

import (
	"container/heap"
)

type RingBufferOrderBook struct {
	Bids      *BidHeap
	Asks      *AskHeap
	InputChan chan Batch
}

type Batch struct {
	Bids []OrderBookEntry
	Asks []OrderBookEntry
}

func NewRingBufferOrderBook(bufferSize int) *RingBufferOrderBook {
	rb := &RingBufferOrderBook{
		Bids:      &BidHeap{},
		Asks:      &AskHeap{},
		InputChan: make(chan Batch, bufferSize),
	}
	heap.Init(rb.Bids)
	heap.Init(rb.Asks)

	go rb.processLoop()
	return rb
}

func (rb *RingBufferOrderBook) AddOrders(bids, asks []OrderBookEntry) {
	rb.InputChan <- Batch{Bids: bids, Asks: asks}
}

func (rb *RingBufferOrderBook) processLoop() {
	for batch := range rb.InputChan {
		for _, b := range batch.Bids {
			heap.Push(rb.Bids, b)
		}
		for _, a := range batch.Asks {
			heap.Push(rb.Asks, a)
		}
	}
}
