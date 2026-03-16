package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/app"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/config"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/persistence"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		if err := runHealthcheck(); err != nil {
			log.Printf("Healthcheck failed: %v", err)
			os.Exit(1)
		}

		return
	}

	if err := run(); err != nil {
		log.Printf("Application failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration initialization failed: %w", err)
	}

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

	log.Println("Database is healthy.")

	repositories := persistence.NewGormRepositories(db.GetDB())
	server := app.NewHTTPServer(cfg, db, repositories)
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Starting server on %s", cfg.Address())
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
		log.Println("Shutdown signal received.")
	}

	gracefulShutdownContext, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(gracefulShutdownContext); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}

	log.Println("Server shutdown complete.")
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
		log.Printf("Database close failed: %v", err)
	}
}
