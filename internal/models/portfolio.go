package models

import (
	"fmt"
	"sync"
)

type Portfolio struct {
	sync.Mutex
	balance float64
	FeeRate float64
}

func NewPortfolio(startUSD float64, feeRate float64) *Portfolio {
	return &Portfolio{
		balance: startUSD,
		FeeRate: feeRate,
	}
}

func (p *Portfolio) ExecuteTrade(buyPrice, sellPrice, quantity float64) (float64, bool) {
	p.Lock()
	defer p.Unlock()

	buyCost := buyPrice * quantity * (1 + p.FeeRate)
	sellRevenue := sellPrice * quantity * (1 - p.FeeRate)
	profit := sellRevenue - buyCost

	if profit <= 0 {
		return profit, false
	}

	if p.balance < buyCost {
		return 0, false
	}

	p.balance += profit

	return profit, true
}

func (p *Portfolio) GetStatus() string {
	p.Lock()
	defer p.Unlock()
	return fmt.Sprintf("💵 Bankroll: $%.2f", p.balance)
}
