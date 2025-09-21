package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cmdb-lite/backend/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// DB holds the database connection
type DB struct {
	*sqlx.DB
}

// Connect establishes a connection to the database
func Connect(cfg *config.Config) (*DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0) // Connections are reused forever

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return &DB{db}, nil
}

// RunMigrations runs database migrations using Goose
func (db *DB) RunMigrations(migrationsDir string) error {
	// Get the underlying sql.DB
	sqlDB := db.DB.DB

	// Set up Goose
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Run migrations
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// GetMigrationDir returns the path to the migrations directory
func GetMigrationDir() (string, error) {
	// The migrations directory is mounted at /database/migrations in the container
	migrationsDir := "/database/migrations"

	// Check if the directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return "", fmt.Errorf("migrations directory does not exist: %s", migrationsDir)
	}

	// Get absolute path
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path to migrations directory: %w", err)
	}

	return absPath, nil
}

// GetSchemaDir returns the path to the schema directory
func GetSchemaDir() (string, error) {
	schemaDir := "/database/schema"

	// Check if the directory exists
	if _, err := os.Stat(schemaDir); os.IsNotExist(err) {
		return "", fmt.Errorf("schema directory does not exist: %s", schemaDir)
	}

	// Get absolute path
	absPath, err := filepath.Abs(schemaDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path to schema directory: %w", err)
	}

	return absPath, nil
}

// GetSeedsDir returns the path to the seeds directory
func GetSeedsDir() (string, error) {
	seedsDir := "/database/seeds"

	// Check if the directory exists
	if _, err := os.Stat(seedsDir); os.IsNotExist(err) {
		return "", fmt.Errorf("seeds directory does not exist: %s", seedsDir)
	}

	// Get absolute path
	absPath, err := filepath.Abs(seedsDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path to seeds directory: %w", err)
	}

	return absPath, nil
}

// ApplySchema applies the initial database schema
func (db *DB) ApplySchema(schemaDir string) error {
	// For now, we'll just run the initial schema file
	// In a real application, you might want to handle multiple schema files
	schemaFile := filepath.Join(schemaDir, "001_initial_schema.sql")

	// Read the schema file
	schemaSQL, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute the schema
	if _, err := db.Exec(string(schemaSQL)); err != nil {
		return fmt.Errorf("failed to apply schema: %w", err)
	}

	log.Println("Database schema applied successfully")
	return nil
}

// ApplySeeds applies seed data to the database
func (db *DB) ApplySeeds(seedsDir string) error {
	// For now, we'll just run the initial seed file
	// In a real application, you might want to handle multiple seed files
	seedFile := filepath.Join(seedsDir, "001_initial_data.sql")

	// Read the seed file
	seedSQL, err := os.ReadFile(seedFile)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}

	// Execute the seed data
	if _, err := db.Exec(string(seedSQL)); err != nil {
		return fmt.Errorf("failed to apply seed data: %w", err)
	}

	log.Println("Database seed data applied successfully")
	return nil
}
