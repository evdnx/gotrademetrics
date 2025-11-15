# GoTradeMetrics

GoTradeMetrics bundles a comprehensive Prometheus instrumentation layer for algorithmic trading systems. It exposes ready-to-use counters, gauges and histograms that describe every stage of a bot lifecycle: order execution, exchange/API calls, retry budgets, anomaly detection, configuration reloads, and more. Each metric set is registered under a single namespace so the resulting `/metrics` endpoint is immediately scrapeable by Prometheus.

## Highlights

- Battery-included metric groups for trading, exchanges, strategies, retry budgets, health checks, logs, backtesting, websockets, databases, security, market data, and anomaly detection.
- Thread-safe helpers (`RecordTrade`, `RecordRetryAttempt`, etc.) so callers never interact with Prometheus primitives directly.
- Simple HTTP integration via `Metrics.Handler()` or `Metrics.SetupMetricsEndpoint`.
- Resettable registry for test environments and the ability to register custom collectors side-by-side with the built-in ones.

## Installation

```bash
go get github.com/evdnx/gotrademetrics
```

Go 1.22+ is recommended.

## Quick Start

```go
package main

import (
	"log"
	"net/http"

	"github.com/evdnx/gotrademetrics"
)

func main() {
	metrics := metrics.NewMetrics("my_bot")

	mux := http.NewServeMux()
	metrics.SetupMetricsEndpoint(mux)

	// Record trading activity
	metrics.RecordTrade()
	metrics.RecordTradeVolume(1.5)
	metrics.UpdatePortfolioValue(12_500)
	metrics.RecordTradingLatency(0.2)

	// Monitor retry budget
	metrics.RecordRetryRequest()
	metrics.RecordRetryAttempt()
	metrics.UpdateRetryBudget(8)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

Visit `http://localhost:8080/metrics` to inspect the exported series.

## Metric Families At A Glance

| Group         | Notable metrics                                                                     |
|---------------|--------------------------------------------------------------------------------------|
| Trading       | Trades executed, total volume, PnL, drawdown, latency, slippage, open/closed poses. |
| Exchange/API  | Request counts, error reasons, latencies, rate-limits, endpoint dimensioning.       |
| Strategy      | Signals, execution errors, profit factor, win rate, drawdown, expectancy.           |
| Retry Budget  | Requests, retries, successes, budget gauge, failure rate.                           |
| Health/Logs   | Component health scores, log volume, queue depths, Zap & buffer stats.              |
| Backtesting   | Runs, simulation time, profitability, Sharpe, trade counts, win rate.               |
| Runtime Ops   | WebSocket, database, circuit-breaker, security, config reload and anomaly metrics.  |

Every metric is namespaced (`<namespace>_<group>_...`) and documented in `metrics.go`.

## Testing

The project includes unit tests for helper methods. Run them with:

```bash
go test ./...
```
