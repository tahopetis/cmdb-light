package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort   int
	DatabaseHost string
	DatabasePort int
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	JWTSecret    string
	
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