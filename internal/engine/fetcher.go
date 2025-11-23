package engine

import (
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

func FetchQuote(id int, exchange models.Exchange, results chan<- models.Quote, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()

	price, err := exchange.FetchPrice()

	results <- models.Quote{
		Venue: models.Venue{
			Name: exchange.GetName(),
		},
		Symbol:    "BTC-USD",
		Price:     price,
		Latency:   time.Since(start),
		Timestamp: time.Now(),
		Error:     err,
	}
}
