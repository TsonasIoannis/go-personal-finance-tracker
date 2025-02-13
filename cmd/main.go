package main

import (
	"log"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	if err := database.InitDB(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	defer database.CloseDB()

	// Check database connection health
	if err := database.CheckDBConnection(); err != nil {
		log.Println("Database is unavailable:", err)
	} else {
		log.Println("Database is healthy!")
	}

	// Create a new Gin router
	r := gin.Default()

	// Register health & readiness routes
	r.GET("/health", handlers.HealthCheckHandler)
	r.GET("/ready", handlers.ReadinessCheckHandler)

	// Default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the Personal Finance Tracker API!"})
	})

	// Start the server
	log.Println("Starting server on :8080")
	r.Run(":8080")
}
