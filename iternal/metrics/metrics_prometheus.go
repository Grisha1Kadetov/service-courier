package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	RateLimitExceededTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of rate limit exceeded",
		},
		[]string{"path", "method"},
	)

	GatewayRetriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_retries_total",
			Help: "Total number of gateway retries",
		},
		[]string{"method", "code"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RateLimitExceededTotal)
	prometheus.MustRegister(GatewayRetriesTotal)
}
