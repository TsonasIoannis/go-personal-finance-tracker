package main

import (
	"log"
	"os"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/auth"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/controllers"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/handlers"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/middleware"
	repositories "github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories/gorm"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/routes"
	services "github.com/TsonasIoannis/go-personal-finance-tracker/internal/services/default"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize a new PostgresDatabase instance
	db := database.NewPostgresDatabase()

	// Connect to the database
	if err := db.Connect(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Ensure database is closed when the program exits
	defer db.Close()

	// Check database connection health
	if err := db.CheckConnection(); err != nil {
		log.Println("Database is unavailable:", err)
	} else {
		log.Println("Database is healthy!")
	}

	// Create a new Gin router
	r := gin.Default()
	tokenManager := auth.NewJWTManager(os.Getenv("JWT_SECRET"), 24*time.Hour)
	authMiddleware := middleware.AuthMiddleware(tokenManager)

	gormDB := db.GetDB()
	// Repositories
	userRepo := repositories.NewUserRepository(gormDB)
	transactionRepo := repositories.NewTransactionRepository(gormDB)
	budgetRepo := repositories.NewGormBudgetRepository(gormDB)

	// Services
	userService := services.NewUserService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo, budgetRepo)
	budgetService := services.NewBudgetService(budgetRepo)

	// Controllers
	userController := controllers.NewUserController(userService, tokenManager)
	transactionController := controllers.NewTransactionController(transactionService)
	budgetController := controllers.NewBudgetController(budgetService)

	// Register API routes
	routes.SetupRoutes(r, authMiddleware, userController, transactionController, budgetController)

	// Register health & readiness routes
	r.GET("/health", handlers.HealthCheckHandler)
	r.GET("/ready", handlers.ReadinessCheckHandler(db))

	// Default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the Personal Finance Tracker API!"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting server on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Println("Failed to start server: ", err)
	}
}
