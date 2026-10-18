package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/version"
)

var (
	// RED metrics for every route. route is the gin route template
	// (/leads/:id) and not the raw path so the label count stays small.
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests handled.",
		},
		[]string{"method", "route", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "route"},
	)

	HTTPInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Requests currently being served.",
		},
	)

	CacheRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lead_cache_requests_total",
			Help: "Redis cache lookups by result (hit, miss, error).",
		},
		[]string{"result"},
	)

	LeadsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "leads_created_total",
			Help: "Leads created successfully.",
		},
	)

	DependencyUp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dependency_up",
			Help: "1 if the dependency answered the last readiness check.",
		},
		[]string{"dependency"},
	)

	buildInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_build_info",
			Help: "Always 1, labels carry the running version.",
		},
		[]string{"version", "commit"},
	)
)

func init() {
	buildInfo.WithLabelValues(version.Version, version.Commit).Set(1)
}
