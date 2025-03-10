package database

import "gorm.io/gorm"

// Database defines an interface for database operations
type Database interface {
	Connect() error
	GetDB() *gorm.DB
	Close() error
	CheckConnection() error
}
