package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Reason defines the reason for a log or API error
type Reason string

const (
	ReasonInternal         Reason = "internal"
	ReasonAPIError         Reason = "api_error"
	ReasonSanitization     Reason = "sanitization"
	ReasonZapSync          Reason = "zap_sync"
	ReasonRateLimit        Reason = "rate_limit"
	ReasonGCPRetry         Reason = "gcp_retry"
	ReasonContextCancelled Reason = "context_cancelled"
	ReasonQueueFull        Reason = "queue_full"
	ReasonInvalidKey       Reason = "invalid_key"
	ReasonNetworkError     Reason = "network_error"
	ReasonCircuitTrip      Reason = "circuit_trip"
	ReasonAnomalyDetected  Reason = "anomaly_detected"
)

// Metrics holds all observability metrics for the trading bot
type Metrics struct {
	Trading        *TradingMetrics
	Exchange       *ExchangeMetrics
	Strategy       *StrategyMetrics
	RetryBudget    *RetryBudgetMetrics
	Health         *HealthMetrics
	Logging        *LoggingMetrics
	Function       *FunctionMetrics
	Config         *ConfigMetrics
	Backtesting    *BacktestingMetrics
	WebSocket      *WebSocketMetrics
	Database       *DatabaseMetrics
	Security       *SecurityMetrics
	Circuit        *CircuitMetrics
	API            *APIMetrics
	Market         *MarketMetrics
	Anomaly        *AnomalyMetrics
	registry       *prometheus.Registry
	namespace      string
	responseWriter http.ResponseWriter
	mu             sync.RWMutex
}

// NewMetrics initializes a new Metrics instance
// It sets up a Prometheus registry and initializes all metric types
func NewMetrics(namespace string) *Metrics {
	if namespace == "" {
		namespace = "cryptobot"
	}

	registry := prometheus.NewRegistry()

	m := &Metrics{
		registry:  registry,
		namespace: namespace,
	}

	m.Trading = newTradingMetrics(namespace, registry)
	m.Exchange = newExchangeMetrics(namespace, registry)
	m.Strategy = newStrategyMetrics(namespace, registry)
	m.RetryBudget = newRetryBudgetMetrics(namespace, registry)
	m.Health = newHealthMetrics(namespace, registry)
	m.Logging = newLoggingMetrics(namespace, registry)
	m.Function = newFunctionMetrics(namespace, registry)
	m.Config = newConfigMetrics(namespace, registry)
	m.Backtesting = newBacktestingMetrics(namespace, registry)
	m.WebSocket = newWebSocketMetrics(namespace, registry)
	m.Database = newDatabaseMetrics(namespace, registry)
	m.Security = newSecurityMetrics(namespace, registry)
	m.Circuit = newCircuitMetrics(namespace, registry)
	m.API = newAPIMetrics(namespace, registry)
	m.Market = newMarketMetrics(namespace, registry)
	m.Anomaly = newAnomalyMetrics(namespace, registry)

	return m
}

// TradingMetrics holds trading-related metrics
type TradingMetrics struct {
	TradesTotal     prometheus.Counter
	PortfolioValue  prometheus.Gauge
	PnL             prometheus.Gauge
	Drawdown        prometheus.Gauge
	PositionSize    prometheus.Gauge
	OrderBookDepth  prometheus.Gauge
	Latency         prometheus.Histogram
	SuccessRate     prometheus.Gauge
	Slippage        prometheus.Histogram
	TradeVolume     prometheus.Counter
	OpenPositions   prometheus.Gauge
	ClosedPositions prometheus.Counter
	OrderExecution  prometheus.Histogram
}

// newTradingMetrics initializes trading metrics
func newTradingMetrics(namespace string, registry *prometheus.Registry) *TradingMetrics {
	m := &TradingMetrics{
		TradesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "trading_trades_total",
			Help:      "Total number of executed trades",
		}),
		PortfolioValue: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "trading_portfolio_value",
			Help:      "Current portfolio value in base currency",
		}),
		PnL: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "trading_pnl",
			Help:      "Current profit and loss",
		}),
		Drawdown: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "trading_drawdown",
			Help:      "Current drawdown percentage",
		}),
		PositionSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "trading_position_size",
			Help:      "Current position size",
		}),
		OrderBookDepth: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "trading_order_book_depth",
			Help:      "Current order book depth",
		}),
		Latency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "trading_latency_seconds",
			Help:      "Trading operation latency in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
		SuccessRate: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "trading_success_rate",
			Help:      "Trade success rate",
		}),
		Slippage: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "trading_slippage",
			Help:      "Trade slippage in percentage",
			Buckets:   prometheus.DefBuckets,
		}),
		TradeVolume: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "trading_volume_total",
			Help:      "Total trading volume in base currency",
		}),
		OpenPositions: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "trading_open_positions",
			Help:      "Number of open trading positions",
		}),
		ClosedPositions: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "trading_closed_positions_total",
			Help:      "Total number of closed trading positions",
		}),
		OrderExecution: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "trading_order_execution_seconds",
			Help:      "Order execution time in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
	}

	registry.MustRegister(
		m.TradesTotal,
		m.PortfolioValue,
		m.PnL,
		m.Drawdown,
		m.PositionSize,
		m.OrderBookDepth,
		m.Latency,
		m.SuccessRate,
		m.Slippage,
		m.TradeVolume,
		m.OpenPositions,
		m.ClosedPositions,
		m.OrderExecution,
	)

	return m
}

// ExchangeMetrics holds exchange-related metrics
type ExchangeMetrics struct {
	RequestsTotal prometheus.CounterVec
	ErrorsTotal   prometheus.CounterVec
	Latency       prometheus.HistogramVec
	RateLimitHits prometheus.CounterVec
}

// newExchangeMetrics initializes exchange metrics
func newExchangeMetrics(namespace string, registry *prometheus.Registry) *ExchangeMetrics {
	m := &ExchangeMetrics{
		RequestsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exchange_requests_total",
			Help:      "Total number of exchange API requests",
		}, []string{"exchange", "endpoint"}),
		ErrorsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exchange_errors_total",
			Help:      "Total number of exchange API errors",
		}, []string{"exchange", "reason"}),
		Latency: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "exchange_latency_seconds",
			Help:      "Exchange API request latency in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"exchange", "endpoint"}),
		RateLimitHits: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exchange_rate_limit_hits_total",
			Help:      "Total number of rate limit hits",
		}, []string{"exchange"}),
	}

	registry.MustRegister(m.RequestsTotal, m.ErrorsTotal, m.Latency, m.RateLimitHits)
	return m
}

// StrategyMetrics holds strategy-related metrics
type StrategyMetrics struct {
	SignalsTotal      prometheus.CounterVec
	TradesTotal       prometheus.CounterVec
	ErrorsTotal       prometheus.CounterVec
	ExecutionErrors   prometheus.CounterVec
	ProfitFactor      prometheus.GaugeVec
	WinRate           prometheus.GaugeVec
	PnL               prometheus.GaugeVec
	MaxDrawdown       prometheus.GaugeVec
	AvgTradeReturn    prometheus.GaugeVec
	AvgTradeHoldTime  prometheus.HistogramVec
	WinLossRatio      prometheus.GaugeVec
	ConsecutiveWins   prometheus.GaugeVec
	ConsecutiveLosses prometheus.GaugeVec
	ProfitPerTrade    prometheus.GaugeVec
	Expectancy        prometheus.GaugeVec
	ReturnVolatility  prometheus.GaugeVec
	TradeFrequency    prometheus.GaugeVec
}

// newStrategyMetrics initializes strategy metrics
func newStrategyMetrics(namespace string, registry *prometheus.Registry) *StrategyMetrics {
	m := &StrategyMetrics{
		SignalsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "strategy_signals_total",
			Help:      "Total number of trading signals generated",
		}, []string{"strategy", "direction"}),
		TradesTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "strategy_trades_total",
			Help:      "Total number of trades executed by strategy",
		}, []string{"strategy"}),
		ErrorsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "strategy_errors_total",
			Help:      "Total number of strategy errors",
		}, []string{"strategy"}),
		ExecutionErrors: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "strategy_execution_errors_total",
			Help:      "Total number of strategy execution errors",
		}, []string{"strategy", "reason"}),
		ProfitFactor: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_profit_factor",
			Help:      "Profit factor of the strategy",
		}, []string{"strategy"}),
		WinRate: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_win_rate",
			Help:      "Win rate of the strategy",
		}, []string{"strategy"}),
		PnL: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_pnl",
			Help:      "Profit and loss of the strategy",
		}, []string{"strategy"}),
		MaxDrawdown: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_max_drawdown",
			Help:      "Maximum drawdown of the strategy",
		}, []string{"strategy"}),
		AvgTradeReturn: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_avg_trade_return",
			Help:      "Average return per trade",
		}, []string{"strategy"}),
		AvgTradeHoldTime: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "strategy_avg_trade_hold_time_seconds",
			Help:      "Average hold time per trade in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"strategy"}),
		WinLossRatio: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_win_loss_ratio",
			Help:      "Win to loss ratio of the strategy",
		}, []string{"strategy"}),
		ConsecutiveWins: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_consecutive_wins",
			Help:      "Number of consecutive winning trades",
		}, []string{"strategy"}),
		ConsecutiveLosses: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_consecutive_losses",
			Help:      "Number of consecutive losing trades",
		}, []string{"strategy"}),
		ProfitPerTrade: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_profit_per_trade",
			Help:      "Average profit per trade",
		}, []string{"strategy"}),
		Expectancy: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_expectancy",
			Help:      "Expected value per trade",
		}, []string{"strategy"}),
		ReturnVolatility: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_return_volatility",
			Help:      "Volatility of trade returns",
		}, []string{"strategy"}),
		TradeFrequency: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "strategy_trade_frequency",
			Help:      "Frequency of trades per unit time",
		}, []string{"strategy"}),
	}

	registry.MustRegister(
		m.SignalsTotal,
		m.TradesTotal,
		m.ErrorsTotal,
		m.ExecutionErrors,
		m.ProfitFactor,
		m.WinRate,
		m.PnL,
		m.MaxDrawdown,
		m.AvgTradeReturn,
		m.AvgTradeHoldTime,
		m.WinLossRatio,
		m.ConsecutiveWins,
		m.ConsecutiveLosses,
		m.ProfitPerTrade,
		m.Expectancy,
		m.ReturnVolatility,
		m.TradeFrequency,
	)

	return m
}

// RetryBudgetMetrics holds retry budget metrics
type RetryBudgetMetrics struct {
	RequestsTotal prometheus.Counter
	RetriesTotal  prometheus.Counter
	SuccessTotal  prometheus.Counter
	BudgetGauge   prometheus.Gauge
	FailureRate   prometheus.Gauge
}

// newRetryBudgetMetrics initializes retry budget metrics
func newRetryBudgetMetrics(namespace string, registry *prometheus.Registry) *RetryBudgetMetrics {
	m := &RetryBudgetMetrics{
		RequestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "retry_budget_requests_total",
			Help:      "Total number of retry budget requests",
		}),
		RetriesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "retry_budget_retries_total",
			Help:      "Total number of retries",
		}),
		SuccessTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "retry_budget_success_total",
			Help:      "Total number of successful retries",
		}),
		BudgetGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "retry_budget_remaining",
			Help:      "Remaining retry budget",
		}),
		FailureRate: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "retry_budget_failure_rate",
			Help:      "Rate of failed retries",
		}),
	}

	registry.MustRegister(m.RequestsTotal, m.RetriesTotal, m.SuccessTotal, m.BudgetGauge, m.FailureRate)
	return m
}

// HealthMetrics holds health-related metrics
type HealthMetrics struct {
	Status prometheus.GaugeVec
}

// newHealthMetrics initializes health metrics
func newHealthMetrics(namespace string, registry *prometheus.Registry) *HealthMetrics {
	m := &HealthMetrics{
		Status: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "health_status",
			Help:      "Health status of components (1=UP, 0.5=DEGRADED, 0=DOWN, -1=UNKNOWN)",
		}, []string{"component"}),
	}

	registry.MustRegister(m.Status)
	return m
}

// LoggingMetrics holds logging-related metrics
type LoggingMetrics struct {
	EntriesTotal     prometheus.CounterVec
	ErrorsTotal      prometheus.CounterVec
	DroppedTotal     prometheus.CounterVec
	QueueSize        prometheus.GaugeVec
	BufferPoolHits   prometheus.Counter
	BufferPoolMisses prometheus.Counter
	ZapSyncErrors    prometheus.Counter
	GCPQueueOverflow prometheus.Counter
}

// newLoggingMetrics initializes logging metrics
func newLoggingMetrics(namespace string, registry *prometheus.Registry) *LoggingMetrics {
	m := &LoggingMetrics{
		EntriesTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "log_entries_total",
			Help:      "Total number of log entries",
		}, []string{"component", "level"}),
		ErrorsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "log_errors_total",
			Help:      "Total number of log errors",
		}, []string{"component", "reason"}),
		DroppedTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "log_dropped_total",
			Help:      "Total number of dropped log entries",
		}, []string{"component", "reason"}),
		QueueSize: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "log_queue_size",
			Help:      "Current size of log queues",
		}, []string{"component", "queue_type"}),
		BufferPoolHits: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "log_buffer_pool_hits_total",
			Help:      "Total number of buffer pool hits",
		}),
		BufferPoolMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "log_buffer_pool_misses_total",
			Help:      "Total number of buffer pool misses",
		}),
		ZapSyncErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "log_zap_sync_errors_total",
			Help:      "Total number of zap sync errors",
		}),
		GCPQueueOverflow: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "log_gcp_queue_overflow_total",
			Help:      "Total number of GCP queue overflows",
		}),
	}

	registry.MustRegister(
		m.EntriesTotal,
		m.ErrorsTotal,
		m.DroppedTotal,
		m.QueueSize,
		m.BufferPoolHits,
		m.BufferPoolMisses,
		m.ZapSyncErrors,
		m.GCPQueueOverflow,
	)

	return m
}

// FunctionMetrics holds function execution metrics
type FunctionMetrics struct {
	Duration    prometheus.HistogramVec
	Errors      prometheus.CounterVec
	Executions  prometheus.CounterVec
	Concurrency prometheus.GaugeVec
}

// newFunctionMetrics initializes function metrics
func newFunctionMetrics(namespace string, registry *prometheus.Registry) *FunctionMetrics {
	m := &FunctionMetrics{
		Duration: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "function_duration_seconds",
			Help:      "Function execution duration in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"function"}),
		Errors: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "function_errors_total",
			Help:      "Total number of function errors",
		}, []string{"function"}),
		Executions: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "function_executions_total",
			Help:      "Total number of function executions",
		}, []string{"function"}),
		Concurrency: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "function_concurrency",
			Help:      "Current number of concurrent function executions",
		}, []string{"function"}),
	}

	registry.MustRegister(m.Duration, m.Errors, m.Executions, m.Concurrency)
	return m
}

// ConfigMetrics holds configuration-related metrics
type ConfigMetrics struct {
	ComponentsRegistered  prometheus.Counter
	ConfigurationErrors   prometheus.Counter
	LoadTotal             prometheus.Counter
	ReloadTotal           prometheus.Counter
	LastReloadTime        prometheus.Gauge
	ValidationErrorsTotal prometheus.Counter
}

// newConfigMetrics initializes configuration metrics
func newConfigMetrics(namespace string, registry *prometheus.Registry) *ConfigMetrics {
	m := &ConfigMetrics{
		ComponentsRegistered: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "config_components_registered_total",
			Help:      "Total number of registered components",
		}),
		ConfigurationErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "config_errors_total",
			Help:      "Total number of configuration errors",
		}),
		LoadTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "config_load_total",
			Help:      "Total number of configuration loads",
		}),
		ReloadTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "config_reload_total",
			Help:      "Total number of configuration reloads",
		}),
		LastReloadTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "config_last_reload_time",
			Help:      "Timestamp of the last configuration reload (Unix seconds)",
		}),
		ValidationErrorsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "config_validation_errors_total",
			Help:      "Total number of configuration validation errors",
		}),
	}

	registry.MustRegister(
		m.ComponentsRegistered,
		m.ConfigurationErrors,
		m.LoadTotal,
		m.ReloadTotal,
		m.LastReloadTime,
		m.ValidationErrorsTotal,
	)

	return m
}

// BacktestingMetrics holds backtesting-related metrics
type BacktestingMetrics struct {
	RunsTotal      prometheus.Counter
	SimulationTime prometheus.Histogram
	ProfitTotal    prometheus.Gauge
	SharpeRatio    prometheus.Gauge
	MaxDrawdown    prometheus.Gauge
	TradeCount     prometheus.Counter
	WinRate        prometheus.Gauge
}

// newBacktestingMetrics initializes backtesting metrics
func newBacktestingMetrics(namespace string, registry *prometheus.Registry) *BacktestingMetrics {
	m := &BacktestingMetrics{
		RunsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "backtesting_runs_total",
			Help:      "Total number of backtesting runs",
		}),
		SimulationTime: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "backtesting_simulation_time_seconds",
			Help:      "Backtesting simulation time in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
		ProfitTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "backtesting_profit_total",
			Help:      "Total profit from backtesting",
		}),
		SharpeRatio: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "backtesting_sharpe_ratio",
			Help:      "Sharpe ratio from backtesting",
		}),
		MaxDrawdown: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "backtesting_max_drawdown",
			Help:      "Maximum drawdown from backtesting",
		}),
		TradeCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "backtesting_trade_count_total",
			Help:      "Total number of trades in backtesting",
		}),
		WinRate: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "backtesting_win_rate",
			Help:      "Win rate from backtesting",
		}),
	}

	registry.MustRegister(
		m.RunsTotal,
		m.SimulationTime,
		m.ProfitTotal,
		m.SharpeRatio,
		m.MaxDrawdown,
		m.TradeCount,
		m.WinRate,
	)

	return m
}

// WebSocketMetrics holds WebSocket-related metrics
type WebSocketMetrics struct {
	ConnectionsTotal  prometheus.Counter
	ActiveConnections prometheus.Gauge
	MessagesSent      prometheus.Counter
	MessagesReceived  prometheus.Counter
	ErrorsTotal       prometheus.CounterVec
	Latency           prometheus.Histogram
}

// newWebSocketMetrics initializes WebSocket metrics
func newWebSocketMetrics(namespace string, registry *prometheus.Registry) *WebSocketMetrics {
	m := &WebSocketMetrics{
		ConnectionsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "websocket_connections_total",
			Help:      "Total number of WebSocket connections",
		}),
		ActiveConnections: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "websocket_active_connections",
			Help:      "Current number of active WebSocket connections",
		}),
		MessagesSent: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "websocket_messages_sent_total",
			Help:      "Total number of WebSocket messages sent",
		}),
		MessagesReceived: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "websocket_messages_received_total",
			Help:      "Total number of WebSocket messages received",
		}),
		ErrorsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "websocket_errors_total",
			Help:      "Total number of WebSocket errors",
		}, []string{"reason"}),
		Latency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "websocket_latency_seconds",
			Help:      "WebSocket message latency in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
	}

	registry.MustRegister(
		m.ConnectionsTotal,
		m.ActiveConnections,
		m.MessagesSent,
		m.MessagesReceived,
		m.ErrorsTotal,
		m.Latency,
	)

	return m
}

// DatabaseMetrics holds database-related metrics
type DatabaseMetrics struct {
	QueriesTotal      prometheus.Counter
	ErrorsTotal       prometheus.CounterVec
	QueryLatency      prometheus.Histogram
	ConnectionsActive prometheus.Gauge
	ConnectionsIdle   prometheus.Gauge
}

// newDatabaseMetrics initializes database metrics
func newDatabaseMetrics(namespace string, registry *prometheus.Registry) *DatabaseMetrics {
	m := &DatabaseMetrics{
		QueriesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "database_queries_total",
			Help:      "Total number of database queries",
		}),
		ErrorsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "database_errors_total",
			Help:      "Total number of database errors",
		}, []string{"reason"}),
		QueryLatency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "database_query_latency_seconds",
			Help:      "Database query latency in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
		ConnectionsActive: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "database_connections_active",
			Help:      "Number of active database connections",
		}),
		ConnectionsIdle: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "database_connections_idle",
			Help:      "Number of idle database connections",
		}),
	}

	registry.MustRegister(
		m.QueriesTotal,
		m.ErrorsTotal,
		m.QueryLatency,
		m.ConnectionsActive,
		m.ConnectionsIdle,
	)

	return m
}

// SecurityMetrics holds security-related metrics
type SecurityMetrics struct {
	KeyGenerationsTotal prometheus.CounterVec
	KeyRevocationsTotal prometheus.CounterVec
	AuthFailures        prometheus.Counter
	SignatureErrors     prometheus.Counter
}

// newSecurityMetrics initializes security metrics
func newSecurityMetrics(namespace string, registry *prometheus.Registry) *SecurityMetrics {
	m := &SecurityMetrics{
		KeyGenerationsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "security_key_generations_total",
			Help:      "Total number of API key generations",
		}, []string{"exchange"}),
		KeyRevocationsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "security_key_revocations_total",
			Help:      "Total number of API key revocations",
		}, []string{"exchange"}),
		AuthFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "security_auth_failures_total",
			Help:      "Total number of authentication failures",
		}),
		SignatureErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "security_signature_errors_total",
			Help:      "Total number of signature generation errors",
		}),
	}

	registry.MustRegister(
		m.KeyGenerationsTotal,
		m.KeyRevocationsTotal,
		m.AuthFailures,
		m.SignatureErrors,
	)

	return m
}

// CircuitMetrics holds circuit breaker-related metrics
type CircuitMetrics struct {
	RequestsTotal    prometheus.CounterVec
	TripsTotal       prometheus.CounterVec
	SuccessesTotal   prometheus.CounterVec
	FailuresTotal    prometheus.CounterVec
	State            prometheus.GaugeVec
	RecoveryAttempts prometheus.CounterVec
	TripDuration     prometheus.HistogramVec
}

// newCircuitMetrics initializes circuit breaker metrics
func newCircuitMetrics(namespace string, registry *prometheus.Registry) *CircuitMetrics {
	m := &CircuitMetrics{
		RequestsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "circuit_breaker_requests_total",
			Help:      "Total number of requests through circuit breaker",
		}, []string{"service"}),
		TripsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "circuit_breaker_trips_total",
			Help:      "Total number of circuit breaker trips",
		}, []string{"service"}),
		SuccessesTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "circuit_breaker_successes_total",
			Help:      "Total number of successful requests through circuit breaker",
		}, []string{"service"}),
		FailuresTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "circuit_breaker_failures_total",
			Help:      "Total number of failed requests through circuit breaker",
		}, []string{"service"}),
		State: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "circuit_breaker_state",
			Help:      "Current state of circuit breaker (0=closed, 1=open, 2=half-open)",
		}, []string{"service"}),
		RecoveryAttempts: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "circuit_breaker_recovery_attempts_total",
			Help:      "Total number of recovery attempts",
		}, []string{"service"}),
		TripDuration: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "circuit_breaker_trip_duration_seconds",
			Help:      "Duration of circuit breaker open state in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"service"}),
	}

	registry.MustRegister(
		m.RequestsTotal,
		m.TripsTotal,
		m.SuccessesTotal,
		m.FailuresTotal,
		m.State,
		m.RecoveryAttempts,
		m.TripDuration,
	)

	return m
}

// APIMetrics holds API call-related metrics
type APIMetrics struct {
	RequestsTotal   prometheus.CounterVec
	ErrorsTotal     prometheus.CounterVec
	Latency         prometheus.HistogramVec
	RequestDuration prometheus.HistogramVec
	RateLimitHits   prometheus.CounterVec
	SuccessRate     prometheus.GaugeVec
}

// newAPIMetrics initializes API call metrics
func newAPIMetrics(namespace string, registry *prometheus.Registry) *APIMetrics {
	m := &APIMetrics{
		RequestsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "api_requests_total",
			Help:      "Total number of API requests",
		}, []string{"service", "endpoint"}),
		ErrorsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "api_errors_total",
			Help:      "Total number of API errors",
		}, []string{"service", "reason"}),
		Latency: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "api_latency_seconds",
			Help:      "API request latency in seconds",
		}, []string{"service", "endpoint"}),
		RequestDuration: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "api_request_duration_seconds",
			Help:      "API request duration in seconds (including retries)",
		}, []string{"service", "endpoint"}),
		RateLimitHits: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "api_rate_limit_hits_total",
			Help:      "Total number of API rate limit hits",
		}, []string{"service"}),
		SuccessRate: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "api_success_rate",
			Help:      "API request success rate",
		}, []string{"service"}),
	}

	registry.MustRegister(
		m.RequestsTotal,
		m.ErrorsTotal,
		m.Latency,
		m.RequestDuration,
		m.RateLimitHits,
		m.SuccessRate,
	)

	return m
}

// MarketMetrics holds market data-related metrics
type MarketMetrics struct {
	PriceUpdatesTotal prometheus.CounterVec
	VolumeTotal       prometheus.CounterVec
	Price             prometheus.GaugeVec
	Volume            prometheus.GaugeVec
	Volatility        prometheus.GaugeVec
	Spread            prometheus.GaugeVec
	Latency           prometheus.HistogramVec
	ErrorsTotal       prometheus.CounterVec
}

// newMarketMetrics initializes market data metrics
func newMarketMetrics(namespace string, registry *prometheus.Registry) *MarketMetrics {
	m := &MarketMetrics{
		PriceUpdatesTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "market_price_updates_total",
			Help:      "Total number of market price updates",
		}, []string{"exchange", "symbol"}),
		VolumeTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "market_volume_total",
			Help:      "Total market volume",
		}, []string{"exchange", "symbol"}),
		Price: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "market_price",
			Help:      "Current market price",
		}, []string{"exchange", "symbol"}),
		Volume: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "market_volume",
			Help:      "Current market volume",
		}, []string{"exchange", "symbol"}),
		Volatility: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "market_volatility",
			Help:      "Market volatility (standard deviation of price)",
		}, []string{"exchange", "symbol"}),
		Spread: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "market_spread",
			Help:      "Market bid-ask spread",
		}, []string{"exchange", "symbol"}),
		Latency: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "market_update_latency_seconds",
			Help:      "Latency of market data updates in seconds",
		}, []string{"exchange", "symbol"}),
		ErrorsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "market_errors_total",
			Help:      "Total number of market data errors",
		}, []string{"exchange", "reason"}),
	}

	registry.MustRegister(
		m.PriceUpdatesTotal,
		m.VolumeTotal,
		m.Price,
		m.Volume,
		m.Volatility,
		m.Spread,
		m.Latency,
		m.ErrorsTotal,
	)

	return m
}

// AnomalyMetrics holds anomaly detection-related metrics
type AnomalyMetrics struct {
	DetectionsTotal   prometheus.CounterVec
	Total             prometheus.CounterVec
	FalsePositives    prometheus.CounterVec
	Scores            prometheus.HistogramVec
	Latency           prometheus.HistogramVec
	ErrorsTotal       prometheus.CounterVec
	DetectionAccuracy prometheus.GaugeVec
}

// newAnomalyMetrics initializes anomaly detection metrics
func newAnomalyMetrics(namespace string, registry *prometheus.Registry) *AnomalyMetrics {
	m := &AnomalyMetrics{
		DetectionsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "anomaly_detections_total",
			Help:      "Total number of anomaly detections",
		}, []string{"type"}),
		Total: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "anomaly_total",
			Help:      "Total number of anomalies processed",
		}, []string{"type"}),
		FalsePositives: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "anomaly_false_positives_total",
			Help:      "Total number of false positive anomaly detections",
		}, []string{"type"}),
		Scores: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "anomaly_scores",
			Help:      "Anomaly detection scores",
			Buckets:   prometheus.LinearBuckets(0, 0.1, 10),
		}, []string{"type"}),
		Latency: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "anomaly_detection_latency_seconds",
			Help:      "Latency of anomaly detection in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"type"}),
		ErrorsTotal: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "anomaly_errors_total",
			Help:      "Total number of anomaly detection errors",
		}, []string{"reason"}),
		DetectionAccuracy: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "anomaly_detection_accuracy",
			Help:      "Accuracy of anomaly detections",
		}, []string{"type"}),
	}

	registry.MustRegister(
		m.DetectionsTotal,
		m.Total,
		m.FalsePositives,
		m.Scores,
		m.Latency,
		m.ErrorsTotal,
		m.DetectionAccuracy,
	)

	return m
}

// RecordTrade increments the executed trades counter
func (m *Metrics) RecordTrade() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.TradesTotal.Inc()
}

// RecordTradeVolume adds the traded volume to the cumulative counter
func (m *Metrics) RecordTradeVolume(volume float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.TradeVolume.Add(volume)
}

// UpdatePortfolioValue sets the latest portfolio valuation
func (m *Metrics) UpdatePortfolioValue(value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.PortfolioValue.Set(value)
}

// UpdatePnL sets the current profit and loss value
func (m *Metrics) UpdatePnL(pnl float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.PnL.Set(pnl)
}

// UpdateDrawdown sets the current drawdown percentage
func (m *Metrics) UpdateDrawdown(drawdown float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.Drawdown.Set(drawdown)
}

// UpdatePositionSize sets the active position size
func (m *Metrics) UpdatePositionSize(size float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.PositionSize.Set(size)
}

// UpdateOrderBookDepth records the most recent observed depth
func (m *Metrics) UpdateOrderBookDepth(depth float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.OrderBookDepth.Set(depth)
}

// RecordTradingLatency observes trading operation latency
func (m *Metrics) RecordTradingLatency(latencySeconds float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.Latency.Observe(latencySeconds)
}

// UpdateTradingSuccessRate sets the trade success rate
func (m *Metrics) UpdateTradingSuccessRate(rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.SuccessRate.Set(rate)
}

// RecordTradingSlippage observes the trade slippage percentage
func (m *Metrics) RecordTradingSlippage(slippage float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.Slippage.Observe(slippage)
}

// UpdateOpenPositions sets the number of open trading positions
func (m *Metrics) UpdateOpenPositions(count float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.OpenPositions.Set(count)
}

// RecordClosedPosition increments the closed positions counter
func (m *Metrics) RecordClosedPosition() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.ClosedPositions.Inc()
}

// RecordOrderExecution observes the order execution duration
func (m *Metrics) RecordOrderExecution(durationSeconds float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Trading.OrderExecution.Observe(durationSeconds)
}

// RecordRetryRequest increments the retry budget request counter
func (m *Metrics) RecordRetryRequest() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RetryBudget.RequestsTotal.Inc()
}

// RecordRetryAttempt increments the retry attempts counter
func (m *Metrics) RecordRetryAttempt() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RetryBudget.RetriesTotal.Inc()
}

// RecordRetrySuccess increments the successful retry counter
func (m *Metrics) RecordRetrySuccess() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RetryBudget.SuccessTotal.Inc()
}

// UpdateRetryBudget sets the remaining retry budget gauge
func (m *Metrics) UpdateRetryBudget(budget float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RetryBudget.BudgetGauge.Set(budget)
}

// UpdateRetryFailureRate sets the retry failure rate gauge
func (m *Metrics) UpdateRetryFailureRate(rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RetryBudget.FailureRate.Set(rate)
}

// RecordFunctionDuration records function execution duration
func (m *Metrics) RecordFunctionDuration(function string, duration float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Function.Duration.WithLabelValues(function).Observe(duration)
}

// RecordFunctionError increments function error counter
func (m *Metrics) RecordFunctionError(function string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Function.Errors.WithLabelValues(function).Inc()
}

// RecordFunctionExecution increments function execution counter
func (m *Metrics) RecordFunctionExecution(function string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Function.Executions.WithLabelValues(function).Inc()
}

// IncrementFunctionConcurrency increments function concurrency gauge
func (m *Metrics) IncrementFunctionConcurrency(function string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Function.Concurrency.WithLabelValues(function).Inc()
}

// DecrementFunctionConcurrency decrements function concurrency gauge
func (m *Metrics) DecrementFunctionConcurrency(function string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Function.Concurrency.WithLabelValues(function).Dec()
}

// UpdateHealthStatus updates component health status
func (m *Metrics) UpdateHealthStatus(component string, status float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Health.Status.WithLabelValues(component).Set(status)
}

// RecordLogEntry records a log entry
func (m *Metrics) RecordLogEntry(component, level string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logging.EntriesTotal.WithLabelValues(component, level).Inc()
}

// RecordLogError records a log error with reason
func (m *Metrics) RecordLogError(component string, reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logging.ErrorsTotal.WithLabelValues(component, string(reason)).Inc()
}

// RecordLogDropped records a dropped log entry
func (m *Metrics) RecordLogDropped(component string, reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logging.DroppedTotal.WithLabelValues(component, string(reason)).Inc()
}

// RecordQueueSize records log queue sizes
func (m *Metrics) RecordQueueSize(component string, criticalQueueSize, nonCriticalQueueSize int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logging.QueueSize.WithLabelValues(component, "critical").Set(float64(criticalQueueSize))
	m.Logging.QueueSize.WithLabelValues(component, "non_critical").Set(float64(nonCriticalQueueSize))
}

// RecordBufferPoolHit records a buffer pool hit
func (m *Metrics) RecordBufferPoolHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logging.BufferPoolHits.Inc()
}

// RecordBufferPoolMiss records a buffer pool miss
func (m *Metrics) RecordBufferPoolMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logging.BufferPoolMisses.Inc()
}

// RecordZapSyncError records a zap sync error
func (m *Metrics) RecordZapSyncError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logging.ZapSyncErrors.Inc()
}

// RecordGCPQueueOverflow records a GCP queue overflow
func (m *Metrics) RecordGCPQueueOverflow() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logging.GCPQueueOverflow.Inc()
}

// RecordExchangeRequest records an exchange API request
func (m *Metrics) RecordExchangeRequest(exchange, endpoint string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Exchange.RequestsTotal.WithLabelValues(exchange, endpoint).Inc()
}

// RecordExchangeError records an exchange API error
func (m *Metrics) RecordExchangeError(exchange string, reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Exchange.ErrorsTotal.WithLabelValues(exchange, string(reason)).Inc()
}

// RecordExchangeLatency records exchange API latency
func (m *Metrics) RecordExchangeLatency(exchange, endpoint string, latency float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Exchange.Latency.WithLabelValues(exchange, endpoint).Observe(latency)
}

// RecordRateLimitHit records an exchange rate limit hit
func (m *Metrics) RecordRateLimitHit(exchange string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Exchange.RateLimitHits.WithLabelValues(exchange).Inc()
}

// RecordStrategySignal records a strategy signal
func (m *Metrics) RecordStrategySignal(strategy, direction string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.SignalsTotal.WithLabelValues(strategy, direction).Inc()
}

// RecordStrategyTrade records a strategy trade
func (m *Metrics) RecordStrategyTrade(strategy string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.TradesTotal.WithLabelValues(strategy).Inc()
}

// RecordStrategyError records a strategy error
func (m *Metrics) RecordStrategyError(strategy string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.ErrorsTotal.WithLabelValues(strategy).Inc()
}

// RecordStrategyExecutionError records a strategy execution error
func (m *Metrics) RecordStrategyExecutionError(strategy string, reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.ExecutionErrors.WithLabelValues(strategy, string(reason)).Inc()
}

// UpdateStrategyProfitFactor updates strategy profit factor
func (m *Metrics) UpdateStrategyProfitFactor(strategy string, factor float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.ProfitFactor.WithLabelValues(strategy).Set(factor)
}

// UpdateStrategyWinRate updates strategy win rate
func (m *Metrics) UpdateStrategyWinRate(strategy string, rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.WinRate.WithLabelValues(strategy).Set(rate)
}

// UpdateStrategyPnL updates strategy profit and loss
func (m *Metrics) UpdateStrategyPnL(strategy string, pnl float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.PnL.WithLabelValues(strategy).Set(pnl)
}

// UpdateStrategyMaxDrawdown updates strategy maximum drawdown
func (m *Metrics) UpdateStrategyMaxDrawdown(strategy string, drawdown float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.MaxDrawdown.WithLabelValues(strategy).Set(drawdown)
}

// UpdateStrategyAvgTradeReturn updates average trade return
func (m *Metrics) UpdateStrategyAvgTradeReturn(strategy string, returnValue float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.AvgTradeReturn.WithLabelValues(strategy).Set(returnValue)
}

// RecordStrategyAvgTradeHoldTime records average trade hold time
func (m *Metrics) RecordStrategyAvgTradeHoldTime(strategy string, holdTime float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.AvgTradeHoldTime.WithLabelValues(strategy).Observe(holdTime)
}

// UpdateStrategyWinLossRatio updates win/loss ratio
func (m *Metrics) UpdateStrategyWinLossRatio(strategy string, ratio float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.WinLossRatio.WithLabelValues(strategy).Set(ratio)
}

// UpdateStrategyConsecutiveWins updates consecutive wins
func (m *Metrics) UpdateStrategyConsecutiveWins(strategy string, count float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.ConsecutiveWins.WithLabelValues(strategy).Set(count)
}

// UpdateStrategyConsecutiveLosses updates consecutive losses
func (m *Metrics) UpdateStrategyConsecutiveLosses(strategy string, count float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.ConsecutiveLosses.WithLabelValues(strategy).Set(count)
}

// UpdateStrategyProfitPerTrade updates profit per trade
func (m *Metrics) UpdateStrategyProfitPerTrade(strategy string, profit float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.ProfitPerTrade.WithLabelValues(strategy).Set(profit)
}

// UpdateStrategyExpectancy updates trade expectancy
func (m *Metrics) UpdateStrategyExpectancy(strategy string, expectancy float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.Expectancy.WithLabelValues(strategy).Set(expectancy)
}

// UpdateStrategyReturnVolatility updates return volatility
func (m *Metrics) UpdateStrategyReturnVolatility(strategy string, volatility float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.ReturnVolatility.WithLabelValues(strategy).Set(volatility)
}

// UpdateStrategyTradeFrequency updates trade frequency
func (m *Metrics) UpdateStrategyTradeFrequency(strategy string, frequency float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Strategy.TradeFrequency.WithLabelValues(strategy).Set(frequency)
}

// RecordBacktestingRun records a backtesting run
func (m *Metrics) RecordBacktestingRun() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Backtesting.RunsTotal.Inc()
}

// RecordBacktestingSimulationTime records backtesting simulation time
func (m *Metrics) RecordBacktestingSimulationTime(duration float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Backtesting.SimulationTime.Observe(duration)
}

// UpdateBacktestingProfit updates backtesting profit
func (m *Metrics) UpdateBacktestingProfit(profit float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Backtesting.ProfitTotal.Set(profit)
}

// UpdateBacktestingSharpeRatio updates backtesting Sharpe ratio
func (m *Metrics) UpdateBacktestingSharpeRatio(ratio float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Backtesting.SharpeRatio.Set(ratio)
}

// UpdateBacktestingMaxDrawdown updates backtesting max drawdown
func (m *Metrics) UpdateBacktestingMaxDrawdown(drawdown float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Backtesting.MaxDrawdown.Set(drawdown)
}

// RecordBacktestingTrade records a backtesting trade
func (m *Metrics) RecordBacktestingTrade() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Backtesting.TradeCount.Inc()
}

// UpdateBacktestingWinRate updates backtesting win rate
func (m *Metrics) UpdateBacktestingWinRate(rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Backtesting.WinRate.Set(rate)
}

// RecordWebSocketConnection records a WebSocket connection
func (m *Metrics) RecordWebSocketConnection() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.WebSocket.ConnectionsTotal.Inc()
}

// UpdateWebSocketActiveConnections updates active WebSocket connections
func (m *Metrics) UpdateWebSocketActiveConnections(count float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.WebSocket.ActiveConnections.Set(count)
}

// RecordWebSocketMessageSent records a sent WebSocket message
func (m *Metrics) RecordWebSocketMessageSent() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.WebSocket.MessagesSent.Inc()
}

// RecordWebSocketMessageReceived records a received WebSocket message
func (m *Metrics) RecordWebSocketMessageReceived() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.WebSocket.MessagesReceived.Inc()
}

// RecordWebSocketError records a WebSocket error
func (m *Metrics) RecordWebSocketError(reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.WebSocket.ErrorsTotal.WithLabelValues(string(reason)).Inc()
}

// RecordWebSocketLatency records WebSocket message latency
func (m *Metrics) RecordWebSocketLatency(latency float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.WebSocket.Latency.Observe(latency)
}

// RecordDatabaseQuery records a database query
func (m *Metrics) RecordDatabaseQuery() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Database.QueriesTotal.Inc()
}

// RecordDatabaseError records a database error
func (m *Metrics) RecordDatabaseError(reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Database.ErrorsTotal.WithLabelValues(string(reason)).Inc()
}

// RecordDatabaseQueryLatency records database query latency
func (m *Metrics) RecordDatabaseQueryLatency(latency float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Database.QueryLatency.Observe(latency)
}

// UpdateDatabaseConnections updates database connection counts
func (m *Metrics) UpdateDatabaseConnections(active, idle float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Database.ConnectionsActive.Set(active)
	m.Database.ConnectionsIdle.Set(idle)
}

// RecordKeyGeneration records an API key generation
func (m *Metrics) RecordKeyGeneration(exchange string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Security.KeyGenerationsTotal.WithLabelValues(exchange).Inc()
}

// RecordKeyRevocation records an API key revocation
func (m *Metrics) RecordKeyRevocation(exchange string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Security.KeyRevocationsTotal.WithLabelValues(exchange).Inc()
}

// RecordAuthFailure records an authentication failure
func (m *Metrics) RecordAuthFailure() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Security.AuthFailures.Inc()
}

// RecordSignatureError records a signature generation error
func (m *Metrics) RecordSignatureError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Security.SignatureErrors.Inc()
}

// RecordConfigurationError records a configuration error
func (m *Metrics) RecordConfigurationError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Config.ConfigurationErrors.Inc()
}

// RecordConfigLoad records a configuration load
func (m *Metrics) RecordConfigLoad() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Config.LoadTotal.Inc()
}

// RecordConfigReload records a configuration reload
func (m *Metrics) RecordConfigReload() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Config.ReloadTotal.Inc()
}

// UpdateConfigLastReloadTime updates the last reload timestamp
func (m *Metrics) UpdateConfigLastReloadTime(timestamp float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Config.LastReloadTime.Set(timestamp)
}

// RecordConfigValidationError records a configuration validation error
func (m *Metrics) RecordConfigValidationError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Config.ValidationErrorsTotal.Inc()
}

// RecordCircuitRequest records a circuit breaker request
func (m *Metrics) RecordCircuitRequest(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Circuit.RequestsTotal.WithLabelValues(service).Inc()
}

// RecordCircuitTrip records a circuit breaker trip
func (m *Metrics) RecordCircuitTrip(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Circuit.TripsTotal.WithLabelValues(service).Inc()
}

// RecordCircuitSuccess records a successful request through circuit breaker
func (m *Metrics) RecordCircuitSuccess(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Circuit.SuccessesTotal.WithLabelValues(service).Inc()
}

// RecordCircuitFailure records a failed request through circuit breaker
func (m *Metrics) RecordCircuitFailure(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Circuit.FailuresTotal.WithLabelValues(service).Inc()
}

// UpdateCircuitState updates circuit breaker state (0=closed, 1=open, 2=half-open)
func (m *Metrics) UpdateCircuitState(service string, state float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Circuit.State.WithLabelValues(service).Set(state)
}

// RecordCircuitRecoveryAttempt records a circuit breaker recovery attempt
func (m *Metrics) RecordCircuitRecoveryAttempt(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Circuit.RecoveryAttempts.WithLabelValues(service).Inc()
}

// RecordCircuitTripDuration records circuit breaker open state duration
func (m *Metrics) RecordCircuitTripDuration(service string, duration float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Circuit.TripDuration.WithLabelValues(service).Observe(duration)
}

// RecordAPIRequest records an API request
func (m *Metrics) RecordAPIRequest(service, endpoint string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.API.RequestsTotal.WithLabelValues(service, endpoint).Inc()
}

// RecordAPIError records an API error
func (m *Metrics) RecordAPIError(service string, reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.API.ErrorsTotal.WithLabelValues(service, string(reason)).Inc()
}

// RecordAPILatency records API request latency
func (m *Metrics) RecordAPILatency(service, endpoint string, latency float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.API.Latency.WithLabelValues(service, endpoint).Observe(latency)
}

// RecordAPIRequestDuration records API request duration (including retries)
func (m *Metrics) RecordAPIRequestDuration(service, endpoint string, duration float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.API.RequestDuration.WithLabelValues(service, endpoint).Observe(duration)
}

// RecordAPIRateLimitHit records an API rate limit hit
func (m *Metrics) RecordAPIRateLimitHit(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.API.RateLimitHits.WithLabelValues(service).Inc()
}

// UpdateAPISuccessRate updates API success rate
func (m *Metrics) UpdateAPISuccessRate(service string, rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.API.SuccessRate.WithLabelValues(service).Set(rate)
}

// RecordMarketPriceUpdate records a market price update
func (m *Metrics) RecordMarketPriceUpdate(exchange, symbol string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Market.PriceUpdatesTotal.WithLabelValues(exchange, symbol).Inc()
}

// RecordMarketVolume records market volume
func (m *Metrics) RecordMarketVolume(exchange, symbol string, volume float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Market.VolumeTotal.WithLabelValues(exchange, symbol).Add(volume)
}

// UpdateMarketPrice updates current market price
func (m *Metrics) UpdateMarketPrice(exchange, symbol string, price float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Market.Price.WithLabelValues(exchange, symbol).Set(price)
}

// UpdateMarketVolume updates current market volume
func (m *Metrics) UpdateMarketVolume(exchange, symbol string, volume float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Market.Volume.WithLabelValues(exchange, symbol).Set(volume)
}

// UpdateMarketVolatility updates market volatility
func (m *Metrics) UpdateMarketVolatility(exchange, symbol string, volatility float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Market.Volatility.WithLabelValues(exchange, symbol).Set(volatility)
}

// UpdateMarketSpread updates market bid-ask spread
func (m *Metrics) UpdateMarketSpread(exchange, symbol string, spread float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Market.Spread.WithLabelValues(exchange, symbol).Set(spread)
}

// RecordMarketUpdateLatency records market data update latency
func (m *Metrics) RecordMarketUpdateLatency(exchange, symbol string, latency float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Market.Latency.WithLabelValues(exchange, symbol).Observe(latency)
}

// RecordMarketError records a market data error
func (m *Metrics) RecordMarketError(exchange string, reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Market.ErrorsTotal.WithLabelValues(exchange, string(reason)).Inc()
}

// RecordAnomalyDetection records an anomaly detection
func (m *Metrics) RecordAnomalyDetection(anomalyType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Anomaly.DetectionsTotal.WithLabelValues(anomalyType).Inc()
}

// RecordAnomalyTotal records total anomalies processed
func (m *Metrics) RecordAnomalyTotal(anomalyType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Anomaly.Total.WithLabelValues(anomalyType).Inc()
}

// RecordAnomalyFalsePositive records a false positive anomaly
func (m *Metrics) RecordAnomalyFalsePositive(anomalyType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Anomaly.FalsePositives.WithLabelValues(anomalyType).Inc()
}

// RecordAnomalyScore records an anomaly detection score
func (m *Metrics) RecordAnomalyScore(anomalyType string, score float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Anomaly.Scores.WithLabelValues(anomalyType).Observe(score)
}

// RecordAnomalyDetectionLatency records anomaly detection latency
func (m *Metrics) RecordAnomalyDetectionLatency(anomalyType string, latency float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Anomaly.Latency.WithLabelValues(anomalyType).Observe(latency)
}

// RecordAnomalyError records an anomaly detection error
func (m *Metrics) RecordAnomalyError(reason Reason) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Anomaly.ErrorsTotal.WithLabelValues(string(reason)).Inc()
}

// UpdateAnomalyDetectionAccuracy updates anomaly detection accuracy
func (m *Metrics) UpdateAnomalyDetectionAccuracy(anomalyType string, accuracy float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Anomaly.DetectionAccuracy.WithLabelValues(anomalyType).Set(accuracy)
}

// Handler returns the Prometheus HTTP handler
func (m *Metrics) Handler() http.Handler {
	m.mu.Lock()
	defer m.mu.Unlock()
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// InstrumentHandler wraps an HTTP handler with metrics
func (m *Metrics) InstrumentHandler(name string, handler http.Handler) http.Handler {
	m.mu.Lock()
	defer m.mu.Unlock()
	return promhttp.InstrumentHandlerDuration(
		m.Function.Duration.MustCurryWith(prometheus.Labels{"function": name}),
		promhttp.InstrumentHandlerCounter(
			m.Function.Errors.MustCurryWith(prometheus.Labels{"function": name}),
			handler,
		),
	)
}

// NewResponseWriter creates a new response writer for metrics
func (m *Metrics) NewResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	m.mu.Lock()
	defer m.mu.Unlock()
	return &responseWriter{ResponseWriter: w}
}

// responseWriter wraps an HTTP ResponseWriter for metrics
type responseWriter struct {
	http.ResponseWriter
}

// WriteHeader records the status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
}

// SetupMetricsEndpoint sets up the metrics endpoint
func (m *Metrics) SetupMetricsEndpoint(mux *http.ServeMux) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mux.Handle("/metrics", m.Handler())
}

// RegisterCustomCollector registers a custom Prometheus collector
func (m *Metrics) RegisterCustomCollector(collector prometheus.Collector) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registry.MustRegister(collector)
}

// UnregisterCollector unregisters a Prometheus collector
func (m *Metrics) UnregisterCollector(collector prometheus.Collector) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.registry.Unregister(collector)
}

// ResetMetrics resets all metrics to their initial state
func (m *Metrics) ResetMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registry = prometheus.NewRegistry()
	m.Trading = newTradingMetrics(m.namespace, m.registry)
	m.Exchange = newExchangeMetrics(m.namespace, m.registry)
	m.Strategy = newStrategyMetrics(m.namespace, m.registry)
	m.RetryBudget = newRetryBudgetMetrics(m.namespace, m.registry)
	m.Health = newHealthMetrics(m.namespace, m.registry)
	m.Logging = newLoggingMetrics(m.namespace, m.registry)
	m.Function = newFunctionMetrics(m.namespace, m.registry)
	m.Config = newConfigMetrics(m.namespace, m.registry)
	m.Backtesting = newBacktestingMetrics(m.namespace, m.registry)
	m.WebSocket = newWebSocketMetrics(m.namespace, m.registry)
	m.Database = newDatabaseMetrics(m.namespace, m.registry)
	m.Security = newSecurityMetrics(m.namespace, m.registry)
	m.Circuit = newCircuitMetrics(m.namespace, m.registry)
	m.API = newAPIMetrics(m.namespace, m.registry)
	m.Market = newMarketMetrics(m.namespace, m.registry)
	m.Anomaly = newAnomalyMetrics(m.namespace, m.registry)
}
