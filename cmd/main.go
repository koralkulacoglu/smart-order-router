package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/engine"
	"github.com/koralkulacoglu/smart-order-router/internal/models"
	"github.com/koralkulacoglu/smart-order-router/internal/models/exchanges"
)

type Job struct {
	Exchange exchanges.Exchange
	Symbol   string
}

func main() {
	portfolio := models.NewPortfolio(1_000_000, 0.0001) // 0.01% fees

	gob := models.NewGlobalOrderBook()

	jobs := []Job{
		{&exchanges.Coinbase{}, "BTC-USD"},
		{&exchanges.Binance{}, "BTCUSDT"},
		{&exchanges.Kraken{}, "XBTUSD"},
	}

	var wg sync.WaitGroup
	stopMatcher := make(chan bool)

	go engine.RunMatcher(gob, portfolio, stopMatcher)

	fmt.Println("--- Starting Fetchers ---")
	for i, job := range jobs {
		wg.Add(1)
		go engine.FetchOrderBook(i+1, job.Exchange, job.Symbol, gob, &wg)
	}

	wg.Wait()
	time.Sleep(2 * time.Second)
	stopMatcher <- true
}
