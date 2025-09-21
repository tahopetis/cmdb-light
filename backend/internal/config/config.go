package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerPort   int
	DatabaseHost string
	DatabasePort int
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	JWTSecret    string
	
	// Token configuration
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	
	// Logging configuration
	LogLevel     string
	LogFormat    string
	
	// Metrics configuration
	MetricsEnabled bool
	MetricsPath    string
	
	// Tracing configuration
	TracingEnabled  bool
	TracingService  string
	TracingEnv      string
	TracingJaegerURL string
	TracingZipkinURL string
	TracingSamplingRate float64
	
	// CORS configuration
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string
	CORSAllowCredentials bool
}

func Load() *Config {
	cfg := &Config{
		ServerPort:   getEnvAsInt("SERVER_PORT", 8080),
		DatabaseHost: getEnv("DB_HOST", "localhost"),
		DatabasePort: getEnvAsInt("DB_PORT", 5432),
		DatabaseUser: getEnv("DB_USER", "cmdb_user"),
		DatabasePass: getEnv("DB_PASSWORD", "cmdb_password"),
		DatabaseName: getEnv("DB_NAME", "cmdb_lite"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		
		// Token configuration
		AccessTokenDuration:  getEnvAsDuration("ACCESS_TOKEN_DURATION", "15m"),
		RefreshTokenDuration: getEnvAsDuration("REFRESH_TOKEN_DURATION", "168h"), // 7 days
		
		// Logging configuration
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		LogFormat:    getEnv("LOG_FORMAT", "json"),
		
		// Metrics configuration
		MetricsEnabled: getEnvAsBool("METRICS_ENABLED", true),
		MetricsPath:    getEnv("METRICS_PATH", "/metrics"),
		
		// Tracing configuration
		TracingEnabled:       getEnvAsBool("TRACING_ENABLED", false),
		TracingService:       getEnv("TRACING_SERVICE", "cmdb-lite"),
		TracingEnv:           getEnv("TRACING_ENV", "development"),
		TracingJaegerURL:     getEnv("TRACING_JAEGER_URL", ""),
		TracingZipkinURL:     getEnv("TRACING_ZIPKIN_URL", ""),
		TracingSamplingRate:  getEnvAsFloat("TRACING_SAMPLING_RATE", 1.0),
		
		// CORS configuration
		CORSAllowedOrigins:   getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		CORSAllowedMethods:   getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		CORSAllowedHeaders:   getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"}),
		CORSAllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
	}
	
	return cfg
}

// DatabaseURL returns the PostgreSQL connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DatabaseHost, c.DatabasePort, c.DatabaseUser, c.DatabasePass, c.DatabaseName)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	log.Printf("Invalid duration value for %s, using default: %s", key, defaultValue)
	value, _ := time.ParseDuration(defaultValue)
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	
	// Simple split by comma - in production, you might want a more sophisticated parser
	return strings.Split(valueStr, ",")
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	log.Printf("Invalid integer value for %s, using default: %d", key, defaultValue)
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	log.Printf("Invalid boolean value for %s, using default: %t", key, defaultValue)
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	log.Printf("Invalid float value for %s, using default: %f", key, defaultValue)
	return defaultValue
}