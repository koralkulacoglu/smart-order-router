package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
	"github.com/koralkulacoglu/smart-order-router/internal/models/exchanges"
)

func FetchOrderBook(id int, exchange exchanges.Exchange, symbol string, gob *models.GlobalOrderBook, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()
	bids, asks, err := exchange.FetchOrderBook(symbol)
	if err != nil {
		fmt.Printf("[Fetcher #%d] %s Error: %v\n", id, exchange.GetName(), err)
		return
	}

	latency := time.Since(start)

	gob.AddOrders(bids, asks)

	fmt.Printf("[Fetcher #%d] %s fetched %d bids, %d asks in %v\n",
		id, exchange.GetName(), len(bids), len(asks), latency)
}
