package config

import "time"

var (
	// Portfolio Defaults
	StartBankroll = 1_000_000.0
	FeeRate       = 0.0001

	// Engine Limits
	FetchInterval    = 1000 * time.Millisecond
	MinOrderQuantity = 0.0000001
	MaxDataStaleness = 2 * time.Second

	// UI Settings
	DashboardFPS         = 10
	DashboardRefreshRate = time.Second / time.Duration(DashboardFPS)
	DashboardLogLimit    = 15
)
