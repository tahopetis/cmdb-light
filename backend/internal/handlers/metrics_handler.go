package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler handles Prometheus metrics requests
type MetricsHandler struct{}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// ServeHTTP serves Prometheus metrics
func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Use the Prometheus HTTP handler to serve metrics
	promhttp.Handler().ServeHTTP(w, r)
}

// MetricsHandlerFunc returns an http.HandlerFunc for serving metrics
func MetricsHandlerFunc() http.HandlerFunc {
	handler := NewMetricsHandler()
	return handler.ServeHTTP
}