package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/engine"
	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

func main() {
	venues := []models.Venue{
		{Name: "Coinbase"},
		{Name: "Binance"},
		{Name: "Kraken"},
		{Name: "Gemini"},
		{Name: "Bybit"},
		{Name: "OKX"},
		{Name: "KuCoin"},
		{Name: "Bitfinex"},
		{Name: "Huobi"},
		{Name: "Gate.io"},
		{Name: "Bitstamp"},
		{Name: "Crypto.com"},
		{Name: "MEXC"},
		{Name: "Bitget"},
		{Name: "Deribit"},
		{Name: "CME Group"},
		{Name: "LMAX Digital"},
		{Name: "Bullish"},
		{Name: "Bakkt"},
		{Name: "Uniswap"},
	}

	const maxConcurrentWorkers = 4

	workQueue := make(chan models.Venue, len(venues))
	quoteStream := make(chan models.Quote, len(venues))

	var wg sync.WaitGroup

	for i := 1; i <= maxConcurrentWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for venue := range workQueue {
				fetchInnerWg := &sync.WaitGroup{}
				fetchInnerWg.Add(1)
				engine.FetchQuote(workerId, venue, quoteStream, fetchInnerWg)
				fetchInnerWg.Wait()
			}
		}(i)
	}

	systemStart := time.Now()

	for _, venue := range venues {
		workQueue <- venue
	}

	close(workQueue)

	go func() {
		wg.Wait()
		close(quoteStream)
	}()

	validQuotes := 0
	var bestQuote models.Quote
	for quote := range quoteStream {
		if quote.Error != nil {
			fmt.Printf("[%-12s] Failed: %v\n", quote.Venue.Name, quote.Error)
			continue
		}

		fmt.Printf("[%-12s] Price: $%.2f | Latency: %v\n", quote.Venue.Name, quote.Price, quote.Latency)

		if validQuotes == 0 || quote.Price < bestQuote.Price {
			bestQuote = quote
		}
		validQuotes++
	}

	systemLatency := time.Since(systemStart)

	fmt.Println()
	if validQuotes > 0 {
		fmt.Printf("Best Venue:      %s\n", bestQuote.Venue.Name)
		fmt.Printf("Best Price: 	 $%.2f\n", bestQuote.Price)
		fmt.Printf("Routes Scanned:  %d/%d\n", validQuotes, len(venues))
		fmt.Printf("System Latency:  %v\n", systemLatency)
	} else {
		fmt.Println("No valid quotes found.")
	}
}
