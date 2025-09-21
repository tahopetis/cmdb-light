package repositories

import (
	"github.com/jmoiron/sqlx"
)

// DB is a wrapper around sqlx.DB for repository operations
type DB struct {
	DB *sqlx.DB
}

// NewDB creates a new DB wrapper
func NewDB(db *sqlx.DB) *DB {
	return &DB{DB: db}
}
