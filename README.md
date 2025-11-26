# Smart Order Router

A concurrent Smart Order Router (SOR) written in Go.

This application simulates an **arbitrage trading engine** that aggregates live order book data from multiple exchanges (Binance, Coinbase, Kraken) in parallel into a centralized "Global Order Book" to identify and execute profitable spread opportunities in real-time.

<video src="Screen Recording 2025-11-26 002828.mp4" title="Title" autoplay muted loop></video>

## How It Works

The system operates as a continuous high-frequency trading simulation:

![alt text](<How It Works.png>)

### Data Ingestion

The system spawns concurrent Fetcher Workers for every exchange. These workers:

- Query external APIs (Binance, Coinbase, Kraken) in parallel.
- Normalize the JSON responses into a standard OrderBookEntry format.
- Push bids (buy offers) and asks (sell offers) into the Global Order Book.

### Global Order Book

This is a thread-safe memory structure protected by a Mutex. It uses two specialized Heaps to organize the data:

- Max-Heap (Bids): Keeps the highest buy price at the top (O(1) access).
- Min-Heap (Asks): Keeps the lowest sell price at the top (O(1) access).

### Arbitrage Execution

The Matcher Engine runs in a continuous loop separate from the fetchers.

1. Peeking: It looks at the top of both heaps.
2. Spread Detection: It checks if Highest Bid > Lowest Ask.
3. Execution: If a profit spread exists, it "executes" the trade for the maximum possible quantity.
4. Cleanup: It automatically removes filled orders or invalidates quotes older than 1 second to prevent stale trading.

## Local Setup

[Install Go](https://go.dev/doc/install) before running these commands:

```bash
git clone https://github.com/koralkulacoglu/smart-order-router.git
cd smart-order-router
go run cmd/main.go
```

### Starting Configurations

You can modify configs such as starting balance and fee rate in the `internal/config/config.go` file.
