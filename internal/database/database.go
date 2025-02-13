package database

import "database/sql"

// Database defines an interface for database operations
type Database interface {
	Connect() error
	CheckConnection() error
	GetDB() *sql.DB
	Close() error
}
