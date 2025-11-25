package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/engine"
	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

func main() {
	exchanges := []models.Exchange{
		&models.Coinbase{},
		&models.Binance{},
		&models.Kraken{},
	}

	const maxConcurrentWorkers = 4

	workQueue := make(chan models.Exchange, len(exchanges))
	quoteStream := make(chan models.Quote, len(exchanges))

	var wg sync.WaitGroup

	for i := 1; i <= maxConcurrentWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for exchange := range workQueue {
				fetchInnerWg := &sync.WaitGroup{}
				fetchInnerWg.Add(1)
				engine.FetchQuote(workerId, exchange, quoteStream, fetchInnerWg)
				fetchInnerWg.Wait()
			}
		}(i)
	}

	systemStart := time.Now()

	for _, exchange := range exchanges {
		workQueue <- exchange
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
			fmt.Printf("[%-12s] Failed: %v\n", quote.Exchange, quote.Error)
			continue
		}

		fmt.Printf("[%-12s] Price: $%.2f | Latency: %v\n", quote.Exchange, quote.Price, quote.Latency)

		if validQuotes == 0 || quote.Price < bestQuote.Price {
			bestQuote = quote
		}
		validQuotes++
	}

	systemLatency := time.Since(systemStart)

	fmt.Println()
	if validQuotes > 0 {
		fmt.Printf("Best Exchange:   %s\n", bestQuote.Exchange)
		fmt.Printf("Best Price: 	 $%.2f\n", bestQuote.Price)
		fmt.Printf("Routes Scanned:  %d/%d\n", validQuotes, len(exchanges))
		fmt.Printf("System Latency:  %v\n", systemLatency)
	} else {
		fmt.Println("No valid quotes found.")
	}
}
