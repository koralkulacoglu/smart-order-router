package exchanges

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

type Binance struct{}

func (b *Binance) GetName() string { return "Binance" }

func (b *Binance) FetchOrderBook(symbol string) ([]models.OrderBookEntry, []models.OrderBookEntry, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/depth?symbol=%s&limit=5", symbol)
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	var data struct {
		Bids [][]string `json:"bids"`
		Asks [][]string `json:"asks"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, nil, err
	}

	now := time.Now()
	var bids, asks []models.OrderBookEntry

	for _, item := range data.Bids {
		p, _ := strconv.ParseFloat(item[0], 64)
		q, _ := strconv.ParseFloat(item[1], 64)
		bids = append(bids, models.OrderBookEntry{Exchange: "Binance", Symbol: symbol, Price: p, Quantity: q, Timestamp: now})
	}
	for _, item := range data.Asks {
		p, _ := strconv.ParseFloat(item[0], 64)
		q, _ := strconv.ParseFloat(item[1], 64)
		asks = append(asks, models.OrderBookEntry{Exchange: "Binance", Symbol: symbol, Price: p, Quantity: q, Timestamp: now})
	}
	return bids, asks, nil
}
