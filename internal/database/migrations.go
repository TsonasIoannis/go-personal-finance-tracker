package database

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type migration struct {
	version string
	name    string
	up      func(*gorm.DB) error
}

var migrations = []migration{
	{
		version: "0001_create_users",
		name:    "create users table",
		up: func(db *gorm.DB) error {
			statements, err := statementsForDialect(db,
				[]string{
					`CREATE TABLE IF NOT EXISTS users (
						id BIGSERIAL PRIMARY KEY,
						name VARCHAR(100) NOT NULL,
						email TEXT NOT NULL,
						password TEXT NOT NULL,
						created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
					)`,
					`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email)`,
				},
				[]string{
					`CREATE TABLE IF NOT EXISTS users (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT NOT NULL,
						email TEXT NOT NULL,
						password TEXT NOT NULL,
						created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
					)`,
					`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email)`,
				},
			)
			return executeStatements(db, statements, err)
		},
	},
	{
		version: "0002_create_categories",
		name:    "create categories table",
		up: func(db *gorm.DB) error {
			statements, err := statementsForDialect(db,
				[]string{
					`CREATE TABLE IF NOT EXISTS categories (
						id BIGSERIAL PRIMARY KEY,
						name VARCHAR(100) NOT NULL,
						description VARCHAR(255)
					)`,
					`CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name ON categories (name)`,
				},
				[]string{
					`CREATE TABLE IF NOT EXISTS categories (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT NOT NULL,
						description TEXT
					)`,
					`CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_name ON categories (name)`,
				},
			)
			return executeStatements(db, statements, err)
		},
	},
	{
		version: "0003_create_transactions",
		name:    "create transactions table",
		up: func(db *gorm.DB) error {
			statements, err := statementsForDialect(db,
				[]string{
					`CREATE TABLE IF NOT EXISTS transactions (
						id BIGSERIAL PRIMARY KEY,
						user_id BIGINT NOT NULL,
						"type" VARCHAR(10) NOT NULL,
						amount DOUBLE PRECISION NOT NULL,
						category_id BIGINT NOT NULL,
						date TIMESTAMPTZ NOT NULL,
						note VARCHAR(255),
						created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
					)`,
					`CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions (user_id)`,
					`CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions (category_id)`,
				},
				[]string{
					`CREATE TABLE IF NOT EXISTS transactions (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						user_id INTEGER NOT NULL,
						"type" TEXT NOT NULL,
						amount REAL NOT NULL,
						category_id INTEGER NOT NULL,
						date DATETIME NOT NULL,
						note TEXT,
						created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
						updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
					)`,
					`CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions (user_id)`,
					`CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions (category_id)`,
				},
			)
			return executeStatements(db, statements, err)
		},
	},
	{
		version: "0004_create_budgets",
		name:    "create budgets table",
		up: func(db *gorm.DB) error {
			statements, err := statementsForDialect(db,
				[]string{
					`CREATE TABLE IF NOT EXISTS budgets (
						id BIGSERIAL PRIMARY KEY,
						user_id BIGINT NOT NULL,
						category_id BIGINT NOT NULL,
						"limit" DOUBLE PRECISION NOT NULL,
						start_date TIMESTAMPTZ NOT NULL,
						end_date TIMESTAMPTZ NOT NULL
					)`,
					`CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets (user_id)`,
					`CREATE INDEX IF NOT EXISTS idx_budgets_category_id ON budgets (category_id)`,
				},
				[]string{
					`CREATE TABLE IF NOT EXISTS budgets (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						user_id INTEGER NOT NULL,
						category_id INTEGER NOT NULL,
						"limit" REAL NOT NULL,
						start_date DATETIME NOT NULL,
						end_date DATETIME NOT NULL
					)`,
					`CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets (user_id)`,
					`CREATE INDEX IF NOT EXISTS idx_budgets_category_id ON budgets (category_id)`,
				},
			)
			return executeStatements(db, statements, err)
		},
	},
}

func ApplyMigrations(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if err := ensureMigrationTable(db); err != nil {
		return fmt.Errorf("ensure schema_migrations table: %w", err)
	}

	appliedVersions, err := loadAppliedVersions(db)
	if err != nil {
		return fmt.Errorf("load applied migrations: %w", err)
	}

	for _, migration := range migrations {
		if appliedVersions[migration.version] {
			continue
		}

		if err := applyMigration(db, migration); err != nil {
			return fmt.Errorf("apply migration %s: %w", migration.version, err)
		}

		slog.Info("applied migration", "version", migration.version, "name", migration.name)
	}

	return nil
}

func applyMigration(db *gorm.DB, migration migration) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := migration.up(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?, CURRENT_TIMESTAMP)`, migration.version).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func loadAppliedVersions(db *gorm.DB) (map[string]bool, error) {
	type appliedMigration struct {
		Version string
	}

	var rows []appliedMigration
	if err := db.Raw(`SELECT version FROM schema_migrations`).Scan(&rows).Error; err != nil {
		return nil, err
	}

	appliedVersions := make(map[string]bool, len(rows))
	for _, row := range rows {
		appliedVersions[row.Version] = true
	}

	return appliedVersions, nil
}

func ensureMigrationTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

func executeStatements(db *gorm.DB, statements []string, err error) error {
	if err != nil {
		return err
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			return err
		}
	}

	return nil
}

func statementsForDialect(db *gorm.DB, postgresStatements, sqliteStatements []string) ([]string, error) {
	switch db.Name() {
	case "postgres":
		return postgresStatements, nil
	case "sqlite":
		return sqliteStatements, nil
	default:
		return nil, fmt.Errorf("unsupported dialect %q", db.Name())
	}
}
