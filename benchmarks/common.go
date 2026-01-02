package benchmarks

import (
	"time"
)

type OrderBook interface {
	AddOrders(bids, asks []OrderBookEntry)
}

type OrderBookEntry struct {
	Exchange  string
	Symbol    string
	Price     float64
	Quantity  float64
	Timestamp time.Time
}

type BidHeap []OrderBookEntry

type AskHeap []OrderBookEntry

func (h BidHeap) Len() int { return len(h) }

func (h BidHeap) Less(i, j int) bool {
	if h[i].Price != h[j].Price {
		return h[i].Price > h[j].Price
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
