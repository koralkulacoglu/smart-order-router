# Smart Order Router

A concurrent Smart Order Router (SOR) written in Go.

This application simulates an **arbitrage trading engine** that aggregates live order book data from multiple exchanges (Binance, Coinbase, Kraken) in parallel into a centralized "Global Order Book" to identify and execute profitable spread opportunities in real-time.

## How It Works

The system operates as a continuous high-frequency trading simulation:

```mermaid
graph TD
    %% Styling - High Contrast
    classDef api fill:#ffcc80,stroke:#ef6c00,stroke-width:2px,color:#000;
    classDef internal fill:#90caf9,stroke:#1565c0,stroke-width:2px,color:#000;
    classDef storage fill:#ce93d8,stroke:#6a1b9a,stroke-width:2px,color:#000;
    classDef logic fill:#a5d6a7,stroke:#2e7d32,stroke-width:2px,color:#000;

    subgraph External [🌐 External Exchanges]
        direction LR
        Binance[Binance API]:::api
        Coinbase[Coinbase API]:::api
        Kraken[Kraken API]:::api
    end

    subgraph App [⚙️ Smart Order Router]

        subgraph Producers [Data Ingestion]
            F1(Fetcher Worker 1):::internal
            F2(Fetcher Worker 2):::internal
            F3(Fetcher Worker 3):::internal
        end

        subgraph Store [🧠 Memory Model]
            GOB[("Global Order Book<br/>(Mutex Guarded)")]:::storage
            Bids{Max-Heap<br/>Bids}:::storage
            Asks{Min-Heap<br/>Asks}:::storage
        end

        subgraph Consumer [⚡ Execution]
            Matcher((Matcher<br/>Engine)):::logic
        end
    end

    %% Connections
    Binance -->|Fetch Orders| F1
    Coinbase -->|Fetch Orders| F2
    Kraken -->|Fetch Orders| F3

    F1 -->|Push Orders| GOB
    F2 -->|Push Orders| GOB
    F3 -->|Push Orders| GOB

    GOB --- Bids
    GOB --- Asks

    Matcher -->|1. Lock & Peek| GOB
    Matcher -->|2. Detect Spread| GOB
    Matcher -->|3. Pop & Execute| Matcher
    Matcher -.->|4. Push Partial Fills| GOB
```

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

[Install Go](https://go.dev/doc/install) before doing this:

```bash
git clone https://github.com/koralkulacoglu/smart-order-router.git
cd smart-order-router
go run cmd/main.go
```

It should output something like this:

```
--- Starting Fetchers ---
--- Matcher Engine Started ---
--- Starting Fetchers ---
--- Matcher Engine Started ---
[Fetcher #1] Coinbase fetched 5 bids, 5 asks in 120.566843ms
[Fetcher #3] Kraken fetched 5 bids, 5 asks in 154.220826ms
>>> EXECUTE: Buy 0.1498 on Coinbase @ 87404.00 -> Sell on Kraken @ 87408.30 | Profit: $0.6440
>>> EXECUTE: Buy 0.0243 on Coinbase @ 87404.02 -> Sell on Kraken @ 87408.30 | Profit: $0.1039
[Fetcher #2] Binance fetched 5 bids, 5 asks in 196.260413ms
>>> EXECUTE: Buy 0.0011 on Coinbase @ 87405.55 -> Sell on Binance @ 87447.66 | Profit: $0.0482
>>> EXECUTE: Buy 0.2610 on Kraken @ 87408.40 -> Sell on Binance @ 87447.66 | Profit: $10.2469
```
