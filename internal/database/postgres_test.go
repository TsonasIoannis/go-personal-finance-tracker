package database

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

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
		err := pgDB.Connect(sql.Open)
		assert.Error(t, err)
		assert.Equal(t, "DATABASE_URL environment variable is not set", err.Error())
	})

	t.Run("should fail if sql.Open fails", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "invalid-url")

		// Attempt to connect to an invalid database
		pgDB := NewPostgresDatabase()
		err := pgDB.Connect(sql.Open)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open DB connection")
	})
	t.Run("should fail if db.Ping fails in Connect", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/dbname")

		// Enable ping monitoring for sqlmock
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer mockDB.Close()

		// Expect db.Ping() to return an error
		mock.ExpectPing().WillReturnError(errors.New("failed to ping DB"))

		// Mock sql.Open to return mockDB
		mockOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
			return mockDB, nil
		}

		// Create PostgresDatabase instance
		pgDB := NewPostgresDatabase()

		// Call Connect() with mock sql.Open()
		err = pgDB.Connect(mockOpen)
		assert.Error(t, err)
		assert.Equal(t, "failed to ping DB: failed to ping DB", err.Error())

		// Ensure expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("should succeed when connection is established", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/dbname")

		// Enable ping monitoring for sqlmock
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer mockDB.Close()

		// Expect db.Ping() to return nil (successful ping)
		mock.ExpectPing().WillReturnError(nil)

		// Mock sql.Open to return mockDB
		mockOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
			return mockDB, nil
		}

		// Create PostgresDatabase instance
		pgDB := NewPostgresDatabase()

		// Call Connect() with mock sql.Open()
		err = pgDB.Connect(mockOpen)
		assert.NoError(t, err) // ✅ This ensures the success case runs

		// Ensure expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCheckConnection(t *testing.T) {
	t.Run("should fail if db is nil", func(t *testing.T) {
		pgDB := NewPostgresDatabase()
		err := pgDB.CheckConnection()
		assert.Error(t, err)
		assert.Equal(t, "database connection is not initialized", err.Error())
	})

	t.Run("should fail if db.Ping fails", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true)) // Enable ping monitoring
		assert.NoError(t, err)
		defer mockDB.Close()

		mock.ExpectPing().WillReturnError(errors.New("database is unreachable"))

		pgDB := &PostgresDatabase{db: mockDB}

		err = pgDB.CheckConnection()
		assert.Error(t, err)
		assert.Equal(t, "database is unreachable: database is unreachable", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should succeed if db.Ping succeeds", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true)) // Enable ping monitoring
		assert.NoError(t, err)
		defer mockDB.Close()

		mock.ExpectPing().WillReturnError(nil)

		pgDB := &PostgresDatabase{db: mockDB}

		err = pgDB.CheckConnection()
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetDB(t *testing.T) {
	t.Run("should return the correct DB instance", func(t *testing.T) {
		mockDB, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		pgDB := &PostgresDatabase{db: mockDB}
		assert.Equal(t, mockDB, pgDB.GetDB())
	})
}

func TestClose(t *testing.T) {
	t.Run("should fail if db.Close fails", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		mock.ExpectClose().WillReturnError(errors.New("failed to close connection"))

		pgDB := &PostgresDatabase{db: mockDB}

		err = pgDB.Close()
		assert.Error(t, err)
		assert.Equal(t, "failed to close connection", err.Error())

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should succeed if db.Close succeeds", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer mockDB.Close()

		mock.ExpectClose().WillReturnError(nil)

		pgDB := &PostgresDatabase{db: mockDB}

		err = pgDB.Close()
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return nil if db is nil", func(t *testing.T) {
		pgDB := NewPostgresDatabase()

		err := pgDB.Close()
		assert.NoError(t, err)
	})
}
