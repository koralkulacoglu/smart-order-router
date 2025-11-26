package engine

import (
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
	"github.com/koralkulacoglu/smart-order-router/internal/models/exchanges"
	"github.com/koralkulacoglu/smart-order-router/internal/ui"
)

func FetchOrderBook(id int, exchange exchanges.Exchange, symbol string, gob *models.GlobalOrderBook, dash *ui.Dashboard, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()
	bids, asks, err := exchange.FetchOrderBook(symbol)
	if err != nil {
		dash.Log("[Fetcher #%d] %s Error: %v", id, exchange.GetName(), err)
		return
	}

	gob.AddOrders(bids, asks)

	latency := time.Since(start)

	dash.Log("[Fetcher #%d] %s Slow Fetch: %v", id, exchange.GetName(), latency)
}
