package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Exchange interface {
	GetName() string
	FetchPrice(symbol string) (float64, error)
}

type Quote struct {
	Exchange  string
	Symbol    string
	Price     float64
	Latency   time.Duration
	Timestamp time.Time
	Error     error
}

type Coinbase struct{}
type Binance struct{}
type Kraken struct{}

func (c *Coinbase) GetName() string { return "Coinbase" }
func (b *Binance) GetName() string  { return "Binance" }
func (k *Kraken) GetName() string   { return "Kraken" }

var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}

func (c *Coinbase) FetchPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s/spot", symbol)

	res, err := httpClient.Get(url)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API status: %s", res.Status)
	}

	var data struct {
		Data struct {
			Amount string `json:"amount"`
		} `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(data.Data.Amount, 64)
}

func (b *Binance) FetchPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)

	res, err := httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API status: %s", res.Status)
	}

	var data struct {
		Price string `json:"price"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(data.Price, 64)
}

func (k *Kraken) FetchPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.kraken.com/0/public/Ticker?pair=%s", symbol)

	res, err := httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API status: %s", res.Status)
	}

	var data struct {
		Result map[string]struct {
			C []string `json:"c"`
		} `json:"result"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return 0, err
	}

	if ticker, ok := data.Result["XXBTZUSD"]; ok && len(ticker.C) > 0 {
		return strconv.ParseFloat(ticker.C[0], 64)
	}

	return 0, fmt.Errorf("price data not found in response")
}
