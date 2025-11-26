package exchanges

import (
	"net/http"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

type Exchange interface {
	GetName() string
	FetchOrderBook(symbol string) ([]models.OrderBookEntry, []models.OrderBookEntry, error)
}

var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}
