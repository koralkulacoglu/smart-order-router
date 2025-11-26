package models

import (
	"container/heap"
	"sync"
	"time"
)

type OrderBookEntry struct {
	Exchange  string
	Symbol    string
	Price     float64
	Quantity  float64
	Timestamp time.Time
}

type BidHeap []OrderBookEntry

type AskHeap []OrderBookEntry

type GlobalOrderBook struct {
	sync.Mutex
	Bids *BidHeap
	Asks *AskHeap
}

func NewGlobalOrderBook() *GlobalOrderBook {
	bids := &BidHeap{}
	asks := &AskHeap{}
	heap.Init(bids)
	heap.Init(asks)
	return &GlobalOrderBook{
		Bids: bids,
		Asks: asks,
	}
}

func (gob *GlobalOrderBook) AddOrders(bids, asks []OrderBookEntry) {
	gob.Lock()
	defer gob.Unlock()
	for _, b := range bids {
		heap.Push(gob.Bids, b)
	}
	for _, a := range asks {
		heap.Push(gob.Asks, a)
	}
}

func (h BidHeap) Len() int { return len(h) }

func (h BidHeap) Less(i, j int) bool {
	if h[i].Price != h[j].Price {
		return h[i].Price > h[j].Price
	}
	if h[i].Quantity != h[j].Quantity {
		return h[i].Quantity > h[j].Quantity
	}
	return h[i].Timestamp.After(h[j].Timestamp)
}

func (h BidHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *BidHeap) Push(x any) { *h = append(*h, x.(OrderBookEntry)) }

func (h *BidHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h AskHeap) Len() int { return len(h) }

func (h AskHeap) Less(i, j int) bool {
	if h[i].Price != h[j].Price {
		return h[i].Price < h[j].Price
	}
	if h[i].Quantity != h[j].Quantity {
		return h[i].Quantity > h[j].Quantity
	}
	return h[i].Timestamp.After(h[j].Timestamp)
}

func (h AskHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *AskHeap) Push(x any) { *h = append(*h, x.(OrderBookEntry)) }

func (h *AskHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
