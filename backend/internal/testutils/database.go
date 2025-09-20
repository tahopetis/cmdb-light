package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB represents a test database container
type TestDB struct {
	Container testcontainers.Container
	DB        *sql.DB
	DSN       string
}

// NewTestDB creates a new PostgreSQL test container
func NewTestDB(ctx context.Context) (*TestDB, error) {
	// Get the test database image from environment or use default
	pgImage := os.Getenv("TEST_POSTGRES_IMAGE")
	if pgImage == "" {
		pgImage = "postgres:15-alpine"
	}

	// Create a PostgreSQL container
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(pgImage),
		postgres.WithDatabase("test_cmdb"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	// Get the connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &TestDB{
		Container: pgContainer,
		DB:        db,
		DSN:       connStr,
	}, nil
}

// Close terminates the container and closes the database connection
func (tdb *TestDB) Close(ctx context.Context) error {
	var errs []error

	// Close database connection
	if tdb.DB != nil {
		if err := tdb.DB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close database: %w", err))
		}
	}

	// Terminate container
	if tdb.Container != nil {
		if err := tdb.Container.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate container: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred while closing test database: %v", errs)
	}

	return nil
}

// RunMigrations runs database migrations using goose
func (tdb *TestDB) RunMigrations(migrationsDir string) error {
	// This would typically use goose to run migrations
	// For now, we'll just log that migrations would be run here
	log.Printf("Migrations would be run from: %s", migrationsDir)
	return nil
}

// TruncateTables truncates all tables in the test database
func (tdb *TestDB) TruncateTables(tables []string) error {
	tx, err := tdb.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return tx.Commit()
}

// SeedTestData inserts test data into the database
func (tdb *TestDB) SeedTestData() error {
	// This would typically insert test data
	// For now, we'll just log that test data would be inserted
	log.Println("Test data would be inserted here")
	return nil
}