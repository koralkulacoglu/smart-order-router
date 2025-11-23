package models

import "time"

type Venue struct {
	Name string
	URL  string
}

type Quote struct {
	Venue     Venue
	Symbol    string
	Price     float64
	Latency   time.Duration
	Timestamp time.Time
	Error     error
}
