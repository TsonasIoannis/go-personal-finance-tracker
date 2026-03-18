package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/app"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/config"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/observability"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/persistence"
)

// @title Personal Finance Tracker API
// @version 1.0
// @description Personal Finance Tracker API with JWT auth, pagination, filtering, and observability hooks.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide the JWT as `Bearer <token>`.

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		if err := runHealthcheck(); err != nil {
			slog.Error("healthcheck failed", "error", err)
			os.Exit(1)
		}

		return
	}

	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration initialization failed: %w", err)
	}

	shutdownTracing, err := observability.SetupTracing(context.Background(), cfg.Tracing)
	if err != nil {
		return fmt.Errorf("tracing initialization failed: %w", err)
	}
	defer func() {
		if shutdownTracing == nil {
			return
		}
		if err := shutdownTracing(context.Background()); err != nil {
			slog.Error("tracing shutdown failed", "error", err)
		}
	}()

	db := database.NewPostgresDatabase(cfg.DatabaseURL)
	if err := db.Connect(); err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}
	defer closeDatabase(db)

	if err := db.Migrate(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	if err := db.CheckConnection(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	slog.Info("database is healthy")

	repositories := persistence.NewGormRepositories(db.GetDB())
	server := app.NewHTTPServer(cfg, db, repositories)
	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("starting server", "address", cfg.Address())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdownContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server failed: %w", err)
	case <-shutdownContext.Done():
		slog.Info("shutdown signal received")
	}

	gracefulShutdownContext, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(gracefulShutdownContext); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	if err := shutdownTracing(gracefulShutdownContext); err != nil {
		slog.Error("tracing shutdown failed", "error", err)
	}
	shutdownTracing = nil

	slog.Info("server shutdown complete")
	return nil
}

func runHealthcheck() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func closeDatabase(db database.Database) {
	if err := db.Close(); err != nil {
		slog.Error("database close failed", "error", err)
	}
}
