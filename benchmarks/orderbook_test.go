package benchmarks

import (
	"testing"
	"time"
)

var (
	symbols = []string{"BTC", "ETH", "SOL", "ADA", "XRP", "DOT", "LTC", "LINK"}
	batches []Batch
)

func init() {
	for _, sym := range symbols {
		bids := []OrderBookEntry{
			{Symbol: sym, Price: 100.0, Quantity: 1.0, Timestamp: time.Now()},
		}
		asks := []OrderBookEntry{
			{Symbol: sym, Price: 101.0, Quantity: 1.0, Timestamp: time.Now()},
		}
		batches = append(batches, Batch{Bids: bids, Asks: asks})
	}
}

func BenchmarkMutexOrderBook(b *testing.B) {
	mob := NewMutexOrderBook()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			batch := batches[i%len(batches)]
			mob.AddOrders(batch.Bids, batch.Asks)
			i++
		}
	})
}

func BenchmarkRingBufferOrderBook(b *testing.B) {
	rob := NewRingBufferOrderBook(1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			batch := batches[i%len(batches)]
			rob.AddOrders(batch.Bids, batch.Asks)
			i++
		}
	})
}

func BenchmarkShardedOrderBook(b *testing.B) {
	sob := NewShardedOrderBook(32)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			batch := batches[i%len(batches)]
			sob.AddOrders(batch.Bids, batch.Asks)
			i++
		}
	})
}
