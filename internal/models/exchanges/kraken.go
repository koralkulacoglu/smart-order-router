package exchanges

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

type Kraken struct{}

func (k *Kraken) GetName() string { return "Kraken" }

func (k *Kraken) FetchOrderBook(symbol string) ([]models.OrderBookEntry, []models.OrderBookEntry, error) {
	url := fmt.Sprintf("https://api.kraken.com/0/public/Depth?pair=%s&count=5", symbol)
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	var data struct {
		Result map[string]struct {
			Bids [][]interface{} `json:"bids"`
			Asks [][]interface{} `json:"asks"`
		} `json:"result"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, nil, err
	}

	now := time.Now()
	var bids, asks []models.OrderBookEntry

	for _, pairData := range data.Result {
		for _, b := range pairData.Bids {
			p, _ := strconv.ParseFloat(b[0].(string), 64)
			q, _ := strconv.ParseFloat(b[1].(string), 64)
			bids = append(bids, models.OrderBookEntry{Exchange: "Kraken", Symbol: symbol, Price: p, Quantity: q, Timestamp: now})
		}
		for _, a := range pairData.Asks {
			p, _ := strconv.ParseFloat(a[0].(string), 64)
			q, _ := strconv.ParseFloat(a[1].(string), 64)
			asks = append(asks, models.OrderBookEntry{Exchange: "Kraken", Symbol: symbol, Price: p, Quantity: q, Timestamp: now})
		}
	}
	return bids, asks, nil
}
