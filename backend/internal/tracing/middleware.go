package tracing

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// TraceIDHeader is the header used to propagate the trace ID
const TraceIDHeader = "X-Trace-ID"

// HTTPTracingMiddleware creates a middleware that traces HTTP requests
func HTTPTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the incoming context from the request
		ctx := r.Context()

		// Generate a trace ID if not already present
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Add trace ID to the response header
		w.Header().Set(TraceIDHeader, traceID)

		// Extract the parent span context from the request headers
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

		// Create a new span for the request
		tracer := otel.Tracer("http")
		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(r.Method),
				semconv.HTTPURLKey.String(r.URL.String()),
				semconv.HTTPUserAgentKey.String(r.UserAgent()),
				semconv.HTTPSchemeKey.String(r.URL.Scheme),
				semconv.HTTPHostKey.String(r.Host),
				semconv.HTTPTargetKey.String(r.URL.Path),
				attribute.String("trace.id", traceID),
				attribute.String("request.id", traceID),
			),
		)
		defer span.End()

		// Add the span to the context
		ctx = context.WithValue(ctx, "span", span)
		ctx = context.WithValue(ctx, "trace.id", traceID)
		ctx = context.WithValue(ctx, "request.id", traceID)

		// Create a response writer wrapper to capture the status code
		rw := &responseWriter{ResponseWriter: w}

		// Call the next handler with the updated context
		start := time.Now()
		next.ServeHTTP(rw, r.WithContext(ctx))
		duration := time.Since(start)

		// Add additional attributes to the span
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int(rw.statusCode),
			attribute.Float64("http.request.duration_ms", float64(duration.Milliseconds())),
		)

		// Set the span status based on the status code
		if rw.statusCode >= 400 && rw.statusCode < 500 {
			span.SetStatus(codes.Error, "Client error")
		} else if rw.statusCode >= 500 {
			span.SetStatus(codes.Error, "Server error")
		}
	})
}

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the status code if not already set
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

// GetTraceIDFromContext returns the trace ID from the context
func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace.id").(string); ok {
		return traceID
	}
	return ""
}

// GetRequestIDFromContext returns the request ID from the context
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value("request.id").(string); ok {
		return requestID
	}
	return ""
}

// GetSpanFromContext returns the span from the context
func GetSpanFromContext(ctx context.Context) trace.Span {
	if span, ok := ctx.Value("span").(trace.Span); ok {
		return span
	}
	return nil
}