package main

import (
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/config"
	"github.com/koralkulacoglu/smart-order-router/internal/engine"
	"github.com/koralkulacoglu/smart-order-router/internal/models"
	"github.com/koralkulacoglu/smart-order-router/internal/models/exchanges"
	"github.com/koralkulacoglu/smart-order-router/internal/ui"
)

type Job struct {
	Exchange exchanges.Exchange
	Symbol   string
}

func main() {
	portfolio := models.NewPortfolio(config.StartBankroll, config.FeeRate)
	dash := ui.NewDashboard(portfolio)
	gob := models.NewGlobalOrderBook()

	jobs := []Job{
		{&exchanges.Coinbase{}, "BTC-USD"},
		{&exchanges.Binance{}, "BTCUSDT"},
		{&exchanges.Kraken{}, "XBTUSD"},
	}

	stopChan := make(chan bool)

	go dash.Run(stopChan)
	go engine.RunMatcher(gob, portfolio, dash, stopChan)

	fetchTicker := time.NewTicker(config.FetchInterval)
	defer fetchTicker.Stop()

	dash.Log("System initialized. Starting fetch loop...")

	for range fetchTicker.C {
		var wg sync.WaitGroup
		for i, job := range jobs {
			wg.Add(1)
			go engine.FetchOrderBook(i+1, job.Exchange, job.Symbol, gob, dash, &wg)
		}
		wg.Wait()
		dash.Log("tick: market data updated")
	}
}
