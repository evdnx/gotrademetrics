package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestTradingMetricsHelpers(t *testing.T) {
	m := NewMetrics("test")

	m.RecordTrade()
	m.RecordTradeVolume(42.5)
	m.UpdatePortfolioValue(5000)
	m.UpdatePnL(250)
	m.UpdateDrawdown(0.12)
	m.UpdatePositionSize(100)
	m.UpdateOrderBookDepth(25)
	m.UpdateTradingSuccessRate(0.9)
	m.UpdateOpenPositions(3)
	m.RecordClosedPosition()

	if got := testutil.ToFloat64(m.Trading.TradesTotal); got != 1 {
		t.Fatalf("expected 1 recorded trade, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.TradeVolume); got != 42.5 {
		t.Fatalf("expected 42.5 trade volume, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.PortfolioValue); got != 5000 {
		t.Fatalf("expected portfolio value 5000, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.PnL); got != 250 {
		t.Fatalf("expected pnl 250, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.Drawdown); got != 0.12 {
		t.Fatalf("expected drawdown 0.12, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.PositionSize); got != 100 {
		t.Fatalf("expected position size 100, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.OrderBookDepth); got != 25 {
		t.Fatalf("expected order book depth 25, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.SuccessRate); got != 0.9 {
		t.Fatalf("expected success rate 0.9, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.OpenPositions); got != 3 {
		t.Fatalf("expected open positions 3, got %v", got)
	}

	if got := testutil.ToFloat64(m.Trading.ClosedPositions); got != 1 {
		t.Fatalf("expected closed positions 1, got %v", got)
	}
}

func TestRetryBudgetHelpers(t *testing.T) {
	m := NewMetrics("test")

	m.RecordRetryRequest()
	m.RecordRetryAttempt()
	m.RecordRetrySuccess()
	m.UpdateRetryBudget(8)
	m.UpdateRetryFailureRate(0.2)

	if got := testutil.ToFloat64(m.RetryBudget.RequestsTotal); got != 1 {
		t.Fatalf("expected 1 retry request, got %v", got)
	}

	if got := testutil.ToFloat64(m.RetryBudget.RetriesTotal); got != 1 {
		t.Fatalf("expected 1 retry attempt, got %v", got)
	}

	if got := testutil.ToFloat64(m.RetryBudget.SuccessTotal); got != 1 {
		t.Fatalf("expected 1 retry success, got %v", got)
	}

	if got := testutil.ToFloat64(m.RetryBudget.BudgetGauge); got != 8 {
		t.Fatalf("expected retry budget 8, got %v", got)
	}

	if got := testutil.ToFloat64(m.RetryBudget.FailureRate); got != 0.2 {
		t.Fatalf("expected retry failure rate 0.2, got %v", got)
	}
}
