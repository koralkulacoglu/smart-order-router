package engine

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

func getPrice() float64 {
	basePrice := 98000.00
	spread := 0.2
	price := basePrice * (1 + rand.Float64()*spread*2 - spread)

	return price
}

func FetchQuote(id int, venue models.Venue, results chan<- models.Quote, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()

	var price float64
	var err error

	// Simulate 10% chance of network failure
	if rand.Float32() > 0.1 {
		price = getPrice()
	} else {
		err = fmt.Errorf("connection failed")
	}

	results <- models.Quote{
		Venue:     venue,
		Symbol:    "BTC-USD",
		Price:     price,
		Latency:   time.Since(start),
		Timestamp: time.Now(),
		Error:     err,
	}
}
