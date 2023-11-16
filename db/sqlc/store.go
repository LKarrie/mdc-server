package db

import (
	"database/sql"
)

// Store provides all func to execute db queries and transactions
type Store interface {
	// Querier
}

// SQLStore provides all func to execute SQL queries and transactions
type SQLStore struct {
	// *Queries
	db *sql.DB
}

// NewStore create a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		// Queries: New(db),
		db: db,
	}
}
