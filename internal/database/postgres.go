package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDatabase implements the Database interface
type PostgresDatabase struct {
	db            *gorm.DB
	connectionURL string
}

// NewPostgresDatabase initializes a new Postgres database instance
func NewPostgresDatabase(connectionURL string) *PostgresDatabase {
	return &PostgresDatabase{connectionURL: connectionURL}
}

// Connect initializes the DB connection using GORM
func (p *PostgresDatabase) Connect() error {
	if p.connectionURL == "" {
		return fmt.Errorf("database connection URL is not configured")
	}

	var err error
	p.db, err = gorm.Open(postgres.Open(p.connectionURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL logging
	})
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %v", err)
	}

	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get generic DB instance: %v", err)
	}

	// Configure connection pooling
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connection established.")
	return nil
}

// GetDB returns the database instance
func (p *PostgresDatabase) GetDB() *gorm.DB {
	return p.db
}

// Check if DB is reachable
func (p *PostgresDatabase) CheckConnection() error {
	if p.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %v", err)
	}

	// Try pinging the database
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database is unreachable: %v", err)
	}

	log.Println("Database connection is active.")
	return nil
}

// Close closes the database connection
func (p *PostgresDatabase) Close() error {
	if p.db == nil {
		return nil // Gracefully return nil if DB is not initialized
	}

	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %v", err)
	}

	return sqlDB.Close()
}

// Migrate applies the configured schema migrations.
func (p *PostgresDatabase) Migrate() error {
	if p.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if err := ApplyMigrations(p.db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations applied successfully.")
	return nil
}
