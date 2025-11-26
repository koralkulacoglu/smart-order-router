package engine

import (
	"container/heap"
	"fmt"
	"math"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/config"
	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

func RunMatcher(gob *models.GlobalOrderBook, portfolio *models.Portfolio, stopChan <-chan bool) {
	fmt.Println("--- Matcher Engine Started ---")
	fmt.Println(portfolio.GetStatus())

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			fmt.Println("--- Matcher Engine Stopped ---")
			fmt.Println("Final " + portfolio.GetStatus())
			return
		case <-ticker.C:
			gob.Lock()

			if gob.Bids.Len() == 0 || gob.Asks.Len() == 0 {
				gob.Unlock()
				continue
			}

			bestBid := (*gob.Bids)[0]
			bestAsk := (*gob.Asks)[0]

			// Check if bid/ask is outdated
			if time.Since(bestBid.Timestamp) > 1*time.Second {
				heap.Pop(gob.Bids)
				gob.Unlock()
				continue
			}
			if time.Since(bestAsk.Timestamp) > 1*time.Second {
				heap.Pop(gob.Asks)
				gob.Unlock()
				continue
			}

			// Check if there is an arbitrage opportunity
			if bestBid.Price >= bestAsk.Price {
				bidOrder := heap.Pop(gob.Bids).(models.OrderBookEntry)
				askOrder := heap.Pop(gob.Asks).(models.OrderBookEntry)

				quantity := math.Min(bidOrder.Quantity, askOrder.Quantity)
				profit, executed := portfolio.ExecuteTrade(askOrder.Price, bidOrder.Price, quantity)

				if executed {
					fmt.Printf(">>> 🚀 EXECUTE: Buy %.4f on %s @ %.2f -> Sell on %s @ %.2f | Profit: $%.4f\n",
						quantity, askOrder.Exchange, askOrder.Price, bidOrder.Exchange, bidOrder.Price, profit)

					bidOrder.Quantity -= quantity
					askOrder.Quantity -= quantity
				}

				if bidOrder.Quantity > config.MinOrderQuantity {
					heap.Push(gob.Bids, bidOrder)
				}
				if askOrder.Quantity > config.MinOrderQuantity {
					heap.Push(gob.Asks, askOrder)
				}
			}

			gob.Unlock()
		}
	}
}
