package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the Prometheus metrics
type Metrics struct {
	// HTTP metrics
	httpRequestsTotal      *prometheus.CounterVec
	httpRequestDuration    *prometheus.HistogramVec
	httpResponseSize       *prometheus.HistogramVec
	httpRequestsInFlight   prometheus.Gauge

	// Database metrics
	dbQueryTotal           *prometheus.CounterVec
	dbQueryDuration        *prometheus.HistogramVec
	dbConnectionsActive    prometheus.Gauge
	dbConnectionsIdle      prometheus.Gauge
	dbConnectionsTotal     prometheus.Gauge

	// Business metrics
	cisTotal               prometheus.Gauge
	relationshipsTotal     prometheus.Gauge
	usersTotal             prometheus.Gauge
	auditLogsTotal         prometheus.Gauge

	// Authentication metrics
	authRequestsTotal      *prometheus.CounterVec
	authFailuresTotal      *prometheus.CounterVec

	// Error metrics
	errorsTotal            *prometheus.CounterVec

	// System metrics
	systemInfo             *prometheus.GaugeVec
}

// NewMetrics creates and initializes all Prometheus metrics
func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP metrics
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cmdb_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "cmdb_http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		httpResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "cmdb_http_response_size_bytes",
				Help:    "Size of HTTP responses in bytes",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			},
			[]string{"method", "endpoint"},
		),
		httpRequestsInFlight: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cmdb_http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		}),

		// Database metrics
		dbQueryTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cmdb_db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"operation", "table"},
		),
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "cmdb_db_query_duration_seconds",
				Help:    "Duration of database queries in seconds",
				Buckets: []float64{0.001, 0.01, 0.1, 1, 10},
			},
			[]string{"operation", "table"},
		),
		dbConnectionsActive: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cmdb_db_connections_active",
			Help: "Number of active database connections",
		}),
		dbConnectionsIdle: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cmdb_db_connections_idle",
			Help: "Number of idle database connections",
		}),
		dbConnectionsTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cmdb_db_connections_total",
			Help: "Total number of database connections",
		}),

		// Business metrics
		cisTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cmdb_cis_total",
			Help: "Total number of configuration items",
		}),
		relationshipsTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cmdb_relationships_total",
			Help: "Total number of relationships",
		}),
		usersTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cmdb_users_total",
			Help: "Total number of users",
		}),
		auditLogsTotal: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "cmdb_audit_logs_total",
			Help: "Total number of audit logs",
		}),

		// Authentication metrics
		authRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cmdb_auth_requests_total",
				Help: "Total number of authentication requests",
			},
			[]string{"type", "status"},
		),
		authFailuresTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cmdb_auth_failures_total",
				Help: "Total number of authentication failures",
			},
			[]string{"type", "reason"},
		),

		// Error metrics
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cmdb_errors_total",
				Help: "Total number of errors",
			},
			[]string{"type", "component"},
		),

		// System metrics
		systemInfo: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cmdb_system_info",
				Help: "System information",
			},
			[]string{"version", "build", "environment"},
		),
	}
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration, size int64) {
	m.httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
	m.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	m.httpResponseSize.WithLabelValues(method, endpoint).Observe(float64(size))
}

// IncrementHTTPRequestsInFlight increments the in-flight requests gauge
func (m *Metrics) IncrementHTTPRequestsInFlight() {
	m.httpRequestsInFlight.Inc()
}

// DecrementHTTPRequestsInFlight decrements the in-flight requests gauge
func (m *Metrics) DecrementHTTPRequestsInFlight() {
	m.httpRequestsInFlight.Dec()
}

// RecordDBQuery records database query metrics
func (m *Metrics) RecordDBQuery(operation, table string, duration time.Duration) {
	m.dbQueryTotal.WithLabelValues(operation, table).Inc()
	m.dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// UpdateDBConnectionStats updates database connection metrics
func (m *Metrics) UpdateDBConnectionStats(active, idle, total int) {
	m.dbConnectionsActive.Set(float64(active))
	m.dbConnectionsIdle.Set(float64(idle))
	m.dbConnectionsTotal.Set(float64(total))
}

// UpdateBusinessMetrics updates business metrics
func (m *Metrics) UpdateBusinessMetrics(cis, relationships, users, auditLogs int) {
	m.cisTotal.Set(float64(cis))
	m.relationshipsTotal.Set(float64(relationships))
	m.usersTotal.Set(float64(users))
	m.auditLogsTotal.Set(float64(auditLogs))
}

// RecordAuthRequest records authentication request metrics
func (m *Metrics) RecordAuthRequest(authType, status string) {
	m.authRequestsTotal.WithLabelValues(authType, status).Inc()
}

// RecordAuthFailure records authentication failure metrics
func (m *Metrics) RecordAuthFailure(authType, reason string) {
	m.authFailuresTotal.WithLabelValues(authType, reason).Inc()
}

// RecordError records error metrics
func (m *Metrics) RecordError(errorType, component string) {
	m.errorsTotal.WithLabelValues(errorType, component).Inc()
}

// UpdateSystemInfo updates system information metrics
func (m *Metrics) UpdateSystemInfo(version, build, environment string) {
	m.systemInfo.WithLabelValues(version, build, environment).Set(1)
}

// DefaultMetrics returns a default metrics instance
var DefaultMetrics = NewMetrics()