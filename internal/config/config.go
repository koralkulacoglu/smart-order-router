package config

import "time"

var (
	// Portfolio Configs
	StartBankroll = 1_000_000.0
	FeeRate       = 0.00001

	// Engine Limits
	MinOrderQuantity = 0.0000001
	MaxDataStaleness = 1 * time.Second
)
