package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// PostgresDatabase implements the Database interface
type PostgresDatabase struct {
	db *sql.DB
}

// NewPostgresDatabase initializes a new Postgres database instance
func NewPostgresDatabase() *PostgresDatabase {
	return &PostgresDatabase{}
}

// Connect initializes the database connection
func (p *PostgresDatabase) Connect() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	p.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open DB connection: %v", err)
	}

	// Verify the connection is alive
	if err = p.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %v", err)
	}

	log.Println("Database connection established successfully.")
	return nil
}

// CheckConnection checks if the database is available
func (p *PostgresDatabase) CheckConnection() error {
	if p.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if err := p.db.Ping(); err != nil {
		return fmt.Errorf("database is unreachable: %v", err)
	}

	return nil
}

// GetDB returns the database instance
func (p *PostgresDatabase) GetDB() *sql.DB {
	return p.db
}

// Close closes the database connection
func (p *PostgresDatabase) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}
