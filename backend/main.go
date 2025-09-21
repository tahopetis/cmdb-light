package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/config"
	"github.com/cmdb-lite/backend/internal/logging"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/cmdb-lite/backend/internal/router"
	"github.com/cmdb-lite/backend/internal/tracing"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logging.NewLogger("cmdb-lite", cfg.LogLevel)
	logger.Info("Starting CMDB Lite application")

	// Initialize tracing
	tracingConfig := tracing.Config{
		Enabled:     cfg.TracingEnabled,
		ServiceName: cfg.TracingService,
		Environment: cfg.TracingEnv,
		JaegerURL:   cfg.TracingJaegerURL,
		ZipkinURL:   cfg.TracingZipkinURL,
		SamplingRate: cfg.TracingSamplingRate,
	}
	if err := tracing.InitializeDefaultTracerProvider(tracingConfig); err != nil {
		logger.WithError(err).Fatal("Failed to initialize tracing")
	}
	defer tracing.ShutdownTracerProvider(context.Background())

	// Get environment variables
	dbURL := getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/cmdb?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	serverPort := getEnv("SERVER_PORT", "8080")

	// Connect to the database
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Test the database connection
	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("Failed to ping database")
	}

	logger.Info("Successfully connected to the database")

	// Create JWT manager
	jwtManager := auth.NewJWTManager(jwtSecret, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)

	// Create password manager
	passwordManager := auth.NewPasswordManager()

	// Create DB wrapper
	dbWrapper := &repositories.DB{DB: db}

	// Setup the router
	r := router.SetupRouter(dbWrapper, jwtManager, passwordManager, cfg)

	// Start the server
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + serverPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.WithField("port", serverPort).Info("Server starting")
	if err := srv.ListenAndServe(); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}