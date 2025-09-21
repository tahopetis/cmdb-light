package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/config"
	"github.com/cmdb-lite/backend/internal/database"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/cmdb-lite/backend/internal/router"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
}

var startTime = time.Now()

func main() {
	log.Println("Starting CMDB Lite backend server...")

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get migrations directory
	migrationsDir, err := database.GetMigrationDir()
	if err != nil {
		log.Fatalf("Failed to get migrations directory: %v", err)
	}

	// Run database migrations
	if err := db.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize auth managers
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	passwordManager := auth.NewPasswordManager()

	// Create repositories wrapper
	repoDB := repositories.NewDB(db.DB)

	// Setup router with all endpoints
	r := router.SetupRouter(repoDB, jwtManager, passwordManager, cfg)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on port %d", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
