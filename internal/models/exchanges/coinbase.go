package exchanges

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

type Coinbase struct{}

func (c *Coinbase) GetName() string { return "Coinbase" }

func (c *Coinbase) FetchOrderBook(symbol string) ([]models.OrderBookEntry, []models.OrderBookEntry, error) {
	url := fmt.Sprintf("https://api.coinbase.com/api/v3/brokerage/market/product_book?product_id=%s&limit=5", symbol)
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	var data struct {
		Pricebook struct {
			Bids []struct {
				Price string `json:"price"`
				Size  string `json:"size"`
			} `json:"bids"`
			Asks []struct {
				Price string `json:"price"`
				Size  string `json:"size"`
			} `json:"asks"`
		} `json:"pricebook"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, nil, err
	}

	now := time.Now()
	var bids, asks []models.OrderBookEntry

	for _, b := range data.Pricebook.Bids {
		p, _ := strconv.ParseFloat(b.Price, 64)
		q, _ := strconv.ParseFloat(b.Size, 64)
		bids = append(bids, models.OrderBookEntry{Exchange: "Coinbase", Symbol: symbol, Price: p, Quantity: q, Timestamp: now})
	}
	for _, a := range data.Pricebook.Asks {
		p, _ := strconv.ParseFloat(a.Price, 64)
		q, _ := strconv.ParseFloat(a.Size, 64)
		asks = append(asks, models.OrderBookEntry{Exchange: "Coinbase", Symbol: symbol, Price: p, Quantity: q, Timestamp: now})
	}
	return bids, asks, nil
}
