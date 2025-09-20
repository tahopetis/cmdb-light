package testutils

import (
	"log"
	"os"
	"strconv"
)

// TestConfig holds configuration for tests
type TestConfig struct {
	ServerPort      int
	DatabaseURL     string
	DatabaseName    string
	DatabaseUser    string
	DatabasePass    string
	JWTSecret       string
	TestMode        bool
	LogLevel        string
}

// NewTestConfig creates a new test configuration
func NewTestConfig() *TestConfig {
	// Set default values
	cfg := &TestConfig{
		ServerPort:      8081, // Different from production to avoid conflicts
		DatabaseURL:     "localhost:5432",
		DatabaseName:    "test_cmdb",
		DatabaseUser:    "test_user",
		DatabasePass:    "test_password",
		JWTSecret:       "test-jwt-secret-not-for-production",
		TestMode:        true,
		LogLevel:        "error", // Only log errors in tests to reduce noise
	}

	// Override with environment variables if they exist
	if port := os.Getenv("TEST_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.ServerPort = p
		}
	}

	if dbURL := os.Getenv("TEST_DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
	}

	if dbName := os.Getenv("TEST_DATABASE_NAME"); dbName != "" {
		cfg.DatabaseName = dbName
	}

	if dbUser := os.Getenv("TEST_DATABASE_USER"); dbUser != "" {
		cfg.DatabaseUser = dbUser
	}

	if dbPass := os.Getenv("TEST_DATABASE_PASSWORD"); dbPass != "" {
		cfg.DatabasePass = dbPass
	}

	if jwtSecret := os.Getenv("TEST_JWT_SECRET"); jwtSecret != "" {
		cfg.JWTSecret = jwtSecret
	}

	if logLevel := os.Getenv("TEST_LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	return cfg
}

// GetDSN returns the database connection string
func (c *TestConfig) GetDSN() string {
	return "host=" + c.DatabaseURL + " user=" + c.DatabaseUser + " password=" + c.DatabasePass + " dbname=" + c.DatabaseName + " sslmode=disable"
}

// MustGetTestConfig returns the test configuration or panics if it can't be created
func MustGetTestConfig() *TestConfig {
	cfg := NewTestConfig()
	
	// Validate the configuration
	if cfg.DatabaseURL == "" || cfg.DatabaseName == "" || cfg.DatabaseUser == "" || cfg.DatabasePass == "" {
		log.Fatal("Invalid test database configuration")
	}
	
	if cfg.JWTSecret == "" {
		log.Fatal("Test JWT secret cannot be empty")
	}
	
	return cfg
}