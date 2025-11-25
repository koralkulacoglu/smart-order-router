package engine

import (
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

func FetchQuote(id int, exchange models.Exchange, symbol string, results chan<- models.Quote, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()

	price, err := exchange.FetchPrice(symbol)

	results <- models.Quote{
		Exchange:  exchange.GetName(),
		Symbol:    symbol,
		Price:     price,
		Latency:   time.Since(start),
		Timestamp: time.Now(),
		Error:     err,
	}
}
