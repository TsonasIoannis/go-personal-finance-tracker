package database

// Database defines an interface for database operations
type Database interface {
	Connect() error
	Migrate() error
	Close() error
	CheckConnection() error
}
