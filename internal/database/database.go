package database

import "database/sql"

// Database defines an interface for database operations
type Database interface {
	Connect(openDB func(driverName, dataSourceName string) (*sql.DB, error)) error
	CheckConnection() error
	GetDB() *sql.DB
	Close() error
}
