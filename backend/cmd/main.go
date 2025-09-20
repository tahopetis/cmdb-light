package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cmdb-lite/backend/internal/config"
	"github.com/cmdb-lite/backend/internal/database"
)

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
	
	// TODO: Initialize routes
	
	// Basic HTTP server for now
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil, // TODO: Add router
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
	
	// TODO: Graceful shutdown with timeout
}