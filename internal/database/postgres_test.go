package database

import (
	"os"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockGormDB creates an in-memory SQLite instance for testing
func MockGormDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
}

func TestNewPostgresDatabase(t *testing.T) {
	t.Run("should initialize a new instance with nil db", func(t *testing.T) {
		db := NewPostgresDatabase()
		assert.NotNil(t, db)
		assert.Nil(t, db.GetDB()) // Initially, DB should be nil
	})
}

func TestConnect(t *testing.T) {
	t.Run("should fail if DATABASE_URL is not set", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")

		pgDB := NewPostgresDatabase()
		err := pgDB.Connect()
		assert.Error(t, err)
		assert.Equal(t, "DATABASE_URL environment variable is not set", err.Error())
	})

	t.Run("should fail with an invalid DSN", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "invalid-dsn") // This will trigger a connection failure

		pgDB := NewPostgresDatabase()
		err := pgDB.Connect()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to DB")
	})
}

func TestCheckConnection(t *testing.T) {
	t.Run("should fail if db is nil", func(t *testing.T) {
		pgDB := NewPostgresDatabase()
		err := pgDB.CheckConnection()
		assert.Error(t, err)
		assert.Equal(t, "database connection is not initialized", err.Error())
	})

	t.Run("should succeed if db is available", func(t *testing.T) {
		mockDB, err := MockGormDB()
		assert.NoError(t, err)

		pgDB := &PostgresDatabase{db: mockDB}

		err = pgDB.CheckConnection()
		assert.NoError(t, err) // DB is reachable
	})
	t.Run("should fail if Ping() returns an error", func(t *testing.T) {
		// Create an invalid mock DB
		brokenDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

		// Close DB to make Ping() fail
		sqlDB, _ := brokenDB.DB()
		sqlDB.Close() // Now Ping() will fail

		pgDB := &PostgresDatabase{db: brokenDB}

		err := pgDB.CheckConnection()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is unreachable") // Correct failure message
	})
}

func TestGetDB(t *testing.T) {
	t.Run("should return the correct DB instance", func(t *testing.T) {
		mockDB, err := MockGormDB()
		assert.NoError(t, err)

		pgDB := &PostgresDatabase{db: mockDB}
		assert.Equal(t, mockDB, pgDB.GetDB())
	})
}

func TestClose(t *testing.T) {
	t.Run("should succeed if db.Close succeeds", func(t *testing.T) {
		mockDB, err := MockGormDB()
		assert.NoError(t, err)

		pgDB := &PostgresDatabase{db: mockDB}

		err = pgDB.Close()
		assert.NoError(t, err) // ✅ Closing should work
	})

	t.Run("should return nil if db is nil", func(t *testing.T) {
		pgDB := NewPostgresDatabase()

		err := pgDB.Close()
		assert.NoError(t, err)
	})
}

func TestRunMigrations(t *testing.T) {
	t.Run("should successfully apply migrations", func(t *testing.T) {
		// Create an in-memory SQLite database for testing
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		pgDB := &PostgresDatabase{db: mockDB}

		// Call runMigrations()
		err = pgDB.runMigrations()
		assert.NoError(t, err)

		// Verify that tables were created
		assert.True(t, pgDB.db.Migrator().HasTable(&models.User{}), "User table should exist")
		assert.True(t, pgDB.db.Migrator().HasTable(&models.Transaction{}), "Transaction table should exist")
		assert.True(t, pgDB.db.Migrator().HasTable(&models.Category{}), "Category table should exist")
		assert.True(t, pgDB.db.Migrator().HasTable(&models.Budget{}), "Budget table should exist")
	})
}
