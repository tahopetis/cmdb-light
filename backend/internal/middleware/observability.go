package middleware

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/cmdb-lite/backend/internal/logging"
	"github.com/cmdb-lite/backend/internal/metrics"
	"github.com/cmdb-lite/backend/internal/tracing"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
)

// ObservabilityMiddleware creates middleware for logging, metrics, and tracing
func ObservabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID if not present
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = logging.GenerateRequestID()
		}

		// Create context with request ID
		ctx := logging.ContextWithRequestID(r.Context(), requestID)

		// Create logger with request ID
		logger := logging.DefaultLogger().WithRequestID(requestID)
		ctx = logging.ContextWithLogger(ctx, logger)

		// Update request with new context
		r = r.WithContext(ctx)

		// Increment in-flight requests counter
		metrics.DefaultMetrics.IncrementHTTPRequestsInFlight()
		defer metrics.DefaultMetrics.DecrementHTTPRequestsInFlight()

		// Wrap response writer to capture status code and response size
		wrappedWriter := &observabilityResponseWriter{ResponseWriter: w}

		// Apply tracing middleware if enabled
		var tracingHandler http.Handler
		if tracing.DefaultTracerProvider != nil {
			tracingHandler = tracing.HTTPTracingMiddleware(next)
		} else {
			tracingHandler = next
		}

		// Call the next handler
		tracingHandler.ServeHTTP(wrappedWriter, r)

		// Calculate request duration
		duration := time.Since(start)

		// Get route from request
		route := mux.CurrentRoute(r)
		endpoint := "unknown"
		if route != nil {
			endpoint, _ = route.GetPathTemplate()
		}

		// Record metrics
		metrics.DefaultMetrics.RecordHTTPRequest(
			r.Method,
			endpoint,
			wrappedWriter.statusCode,
			duration,
			wrappedWriter.responseSize,
		)

		// Log request
		logger.LogHTTPRequest(
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			r.UserAgent(),
			r.RemoteAddr,
			wrappedWriter.statusCode,
			duration,
			wrappedWriter.responseSize,
		)

		// Add tracing context if enabled
		if tracing.DefaultTracerProvider != nil {
			span := tracing.GetSpanFromContext(r.Context())
			if span != nil {
				span.SetAttributes(
					attribute.String("request_id", requestID),
					attribute.String("endpoint", endpoint),
					attribute.String("method", r.Method),
					attribute.Int("status_code", wrappedWriter.statusCode),
					attribute.Float64("duration_ms", float64(duration.Milliseconds())),
				)
			}
		}
	})
}

// observabilityResponseWriter wraps http.ResponseWriter to capture status code and response size
type observabilityResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int64
}

// WriteHeader captures the status code
func (rw *observabilityResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the response size
func (rw *observabilityResponseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.responseSize += int64(size)
	return size, err
}

// Hijack implements http.Hijacker interface
func (rw *observabilityResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("response writer cannot hijack")
}

// Flush implements http.Flusher interface
func (rw *observabilityResponseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// DatabaseQueryMiddleware creates middleware for database query logging and metrics
func DatabaseQueryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get logger from context
		logger := logging.GetLoggerFromContext(r.Context())

		// Create context with database query hook
		ctx := context.WithValue(r.Context(), "dbQueryHook", func(query string, args []interface{}, duration time.Duration, err error) {
			// Log database query
			logger.LogDBQuery(query, args, duration, err)

			// Record metrics
			// Extract table name from query (simplified)
			table := "unknown"
			if len(query) > 6 {
				// Simple extraction, in a real implementation you would use a SQL parser
				switch query[0:6] {
				case "SELECT":
					// Extract table name after FROM
					fromIndex := findWordIndex(query, "FROM")
					if fromIndex != -1 {
						table = extractTableName(query, fromIndex+4)
					}
				case "INSERT":
					// Extract table name after INTO
					intoIndex := findWordIndex(query, "INTO")
					if intoIndex != -1 {
						table = extractTableName(query, intoIndex+4)
					}
				case "UPDATE":
					// Extract table name after UPDATE
					table = extractTableName(query, 6)
				case "DELETE":
					// Extract table name after FROM
					fromIndex := findWordIndex(query, "FROM")
					if fromIndex != -1 {
						table = extractTableName(query, fromIndex+4)
					}
				}
			}

			// Determine operation type
			operation := "unknown"
			if len(query) > 6 {
				switch query[0:6] {
				case "SELECT":
					operation = "select"
				case "INSERT":
					operation = "insert"
				case "UPDATE":
					operation = "update"
				case "DELETE":
					operation = "delete"
				}
			}

			// Record metrics
			metrics.DefaultMetrics.RecordDBQuery(operation, table, duration)

			// Add tracing context if enabled
			if tracing.DefaultTracerProvider != nil {
				tracing.AddSpanEvent(r.Context(), "database_query",
					attribute.String("query", query),
					attribute.String("table", table),
					attribute.String("operation", operation),
					attribute.Int64("duration_ms", duration.Milliseconds()),
				)

				if err != nil {
					tracing.SetSpanError(r.Context(), err)
				}
			}
		})

		// Update request with new context
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// findWordIndex finds the index of a word in a string
func findWordIndex(s, word string) int {
	for i := 0; i <= len(s)-len(word); i++ {
		if s[i:i+len(word)] == word {
			// Check if it's a whole word
			if (i == 0 || s[i-1] == ' ') && (i+len(word) == len(s) || s[i+len(word)] == ' ') {
				return i
			}
		}
	}
	return -1
}

// extractTableName extracts the table name from a SQL query
func extractTableName(query string, startIndex int) string {
	if startIndex >= len(query) {
		return "unknown"
	}

	// Skip whitespace
	for startIndex < len(query) && query[startIndex] == ' ' {
		startIndex++
	}

	// Extract table name until whitespace or special character
	endIndex := startIndex
	for endIndex < len(query) && query[endIndex] != ' ' && query[endIndex] != ',' && query[endIndex] != ';' && query[endIndex] != '(' {
		endIndex++
	}

	if startIndex >= endIndex {
		return "unknown"
	}

	return query[startIndex:endIndex]
}

// GetDBQueryHook returns the database query hook from context
func GetDBQueryHook(ctx context.Context) func(query string, args []interface{}, duration time.Duration, err error) {
	if hook, ok := ctx.Value("dbQueryHook").(func(query string, args []interface{}, duration time.Duration, err error)); ok {
		return hook
	}
	return nil
}
