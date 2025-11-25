package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/engine"
	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

type Job struct {
	Exchange models.Exchange
	Symbol   string
}

func main() {
	jobs := []Job{
		{&models.Coinbase{}, "BTC-USD"},
		{&models.Binance{}, "BTCUSDT"},
		{&models.Kraken{}, "XBTUSD"},
	}

	const maxConcurrentWorkers = 4

	workQueue := make(chan Job, len(jobs))
	quoteStream := make(chan models.Quote, len(jobs))

	var wg sync.WaitGroup

	for i := 1; i <= maxConcurrentWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for job := range workQueue {
				fetchInnerWg := &sync.WaitGroup{}
				fetchInnerWg.Add(1)
				engine.FetchQuote(workerId, job.Exchange, job.Symbol, quoteStream, fetchInnerWg)
				fetchInnerWg.Wait()
			}
		}(i)
	}

	systemStart := time.Now()

	for _, job := range jobs {
		workQueue <- job
	}

	close(workQueue)

	go func() {
		wg.Wait()
		close(quoteStream)
	}()

	fmt.Println()
	fmt.Println("------------------------------------------------------------------------")
	fmt.Printf("%-12s | %-10s | %-15s | %-12s\n", "EXCHANGE", "SYMBOL", "PRICE", "LATENCY")
	fmt.Println("------------------------------------------------------------------------")

	validQuotes := 0
	var bestQuote models.Quote
	for quote := range quoteStream {
		if quote.Error != nil {
			fmt.Printf("%-12s | %-10s | %-15s | %-12s\n",
				quote.Exchange, quote.Symbol, "FAILED", quote.Latency)
			continue
		}

		fmt.Printf("%-12s | %-10s | $%-14.2f | %-12s\n",
			quote.Exchange, quote.Symbol, quote.Price, quote.Latency)

		if validQuotes == 0 || quote.Price < bestQuote.Price {
			bestQuote = quote
		}
		validQuotes++
	}

	systemLatency := time.Since(systemStart)

	fmt.Println("------------------------------------------------------------------------")
	fmt.Println()

	if validQuotes > 0 {
		fmt.Printf("BEST OFFER:   	%s on %s\n", bestQuote.Symbol, bestQuote.Exchange)
		fmt.Printf("PRICE:          $%.2f\n", bestQuote.Price)
		fmt.Printf("SYSTEM LATENCY: %v\n", systemLatency)
	} else {
		fmt.Println("No valid quotes found.")
	}
}
