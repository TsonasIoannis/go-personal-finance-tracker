package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDatabase implements the Database interface
type PostgresDatabase struct {
	db *gorm.DB
}

// NewPostgresDatabase initializes a new Postgres database instance
func NewPostgresDatabase() *PostgresDatabase {
	return &PostgresDatabase{}
}

// Connect initializes the DB connection using GORM
func (p *PostgresDatabase) Connect() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	p.db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
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

	// Run Migrations
	return p.runMigrations()
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

// runMigrations applies database schema changes
func (p *PostgresDatabase) runMigrations() error {
	err := p.db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.Category{},
		&models.Budget{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	log.Println("Database migrations applied successfully.")
	return nil
}
