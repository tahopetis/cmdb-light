package logging

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "requestID"
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "userID"
)

// Logger is a wrapper around zerolog.Logger with additional functionality
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates a new logger instance
func NewLogger(service string, level string) *Logger {
	// Set log level based on environment
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	// Configure zerolog
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339

	// Create logger
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", service).
		Logger()

	return &Logger{
		logger: logger,
	}
}

// WithRequestID adds request ID to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	newLogger := l.logger.With().Str("request_id", requestID).Logger()
	return &Logger{logger: newLogger}
}

// WithUserID adds user ID to the logger
func (l *Logger) WithUserID(userID string) *Logger {
	newLogger := l.logger.With().Str("user_id", userID).Logger()
	return &Logger{logger: newLogger}
}

// WithField adds a custom field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := l.logger.With().Interface(key, value).Logger()
	return &Logger{logger: newLogger}
}

// WithFields adds multiple custom fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	event := l.logger.With()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	newLogger := event.Logger()
	return &Logger{logger: newLogger}
}

// WithError adds error to the logger
func (l *Logger) WithError(err error) *Logger {
	newLogger := l.logger.With().Err(err).Logger()
	return &Logger{logger: newLogger}
}

// WithCaller adds caller information to the logger
func (l *Logger) WithCaller() *Logger {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}
	newLogger := l.logger.With().Str("file", file).Int("line", line).Logger()
	return &Logger{logger: newLogger}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf logs a debug message with formatting
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof logs an info message with formatting
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf logs a warning message with formatting
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Msgf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf logs an error message with formatting
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Msgf(format, v...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf logs a fatal message with formatting and exits
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Msgf(format, v...)
}

// GetLoggerFromContext returns logger from context
func GetLoggerFromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value("logger").(*Logger); ok {
		return logger
	}
	return DefaultLogger()
}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

// GetRequestIDFromContext returns request ID from context
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// ContextWithRequestID adds request ID to context
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetUserIDFromContext returns user ID from context
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// ContextWithUserID adds user ID to context
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// DefaultLogger returns a default logger instance
func DefaultLogger() *Logger {
	return NewLogger("cmdb-lite", "info")
}

// GenerateRequestID generates a new request ID
func GenerateRequestID() string {
	return uuid.New().String()
}

// LogHTTPRequest logs HTTP request information
func (l *Logger) LogHTTPRequest(method, path, query, userAgent, remoteAddr string, statusCode int, latency time.Duration, contentLength int64) {
	l.logger.Info().
		Str("method", method).
		Str("path", path).
		Str("query", query).
		Str("user_agent", userAgent).
		Str("remote_addr", remoteAddr).
		Int("status_code", statusCode).
		Dur("latency", latency).
		Int64("content_length", contentLength).
		Msg("HTTP request")
}

// LogDBQuery logs database query information
func (l *Logger) LogDBQuery(query string, args []interface{}, duration time.Duration, err error) {
	event := l.logger.Info().
		Str("query", query).
		Interface("args", args).
		Dur("duration", duration)

	if err != nil {
		event = event.Err(err)
	}

	event.Msg("Database query")
}

// LogBusinessEvent logs business events
func (l *Logger) LogBusinessEvent(eventType, action, entityType, entityID string, details map[string]interface{}) {
	event := l.logger.Info().
		Str("event_type", eventType).
		Str("action", action).
		Str("entity_type", entityType).
		Str("entity_id", entityID)

	if details != nil {
		event = event.Interface("details", details)
	}

	event.Msg("Business event")
}

// LogSecurityEvent logs security events
func (l *Logger) LogSecurityEvent(eventType, username, clientIP string, success bool, details map[string]interface{}) {
	event := l.logger.Info().
		Str("event_type", eventType).
		Str("username", username).
		Str("client_ip", clientIP).
		Bool("success", success)

	if details != nil {
		event = event.Interface("details", details)
	}

	event.Msg("Security event")
}
