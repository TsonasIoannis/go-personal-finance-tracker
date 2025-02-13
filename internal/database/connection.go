package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds database connection details
var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %v", err)
	}

	// Verify the connection is alive
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %v", err)
	}

	log.Println("Database connection established successfully.")
	return nil
}

// CheckDBConnection checks if the database is available
func CheckDBConnection() error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database is unreachable: %v", err)
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
