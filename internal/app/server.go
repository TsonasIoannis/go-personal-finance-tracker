package app

import (
	"net/http"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/auth"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/config"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/controllers"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/handlers"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/middleware"
	repositories "github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories/gorm"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/routes"
	services "github.com/TsonasIoannis/go-personal-finance-tracker/internal/services/default"
	"github.com/gin-gonic/gin"
)

func NewHTTPServer(cfg config.Config, db database.Database) *http.Server {
	router := newRouter(cfg, db)

	return &http.Server{
		Addr:              cfg.Address(),
		Handler:           router,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
	}
}

func newRouter(cfg config.Config, db database.Database) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	tokenManager := auth.NewJWTManager(cfg.JWTSecret, cfg.Auth.TokenTTL)
	authMiddleware := middleware.AuthMiddleware(tokenManager)

	gormDB := db.GetDB()

	userRepo := repositories.NewUserRepository(gormDB)
	transactionRepo := repositories.NewTransactionRepository(gormDB)
	budgetRepo := repositories.NewGormBudgetRepository(gormDB)

	userService := services.NewUserService(userRepo)
	transactionService := services.NewTransactionService(transactionRepo, budgetRepo)
	budgetService := services.NewBudgetService(budgetRepo)

	userController := controllers.NewUserController(userService, tokenManager)
	transactionController := controllers.NewTransactionController(transactionService)
	budgetController := controllers.NewBudgetController(budgetService)

	routes.SetupRoutes(router, authMiddleware, userController, transactionController, budgetController)

	router.GET("/health", handlers.HealthCheckHandler)
	router.GET("/ready", handlers.ReadinessCheckHandler(db))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Personal Finance Tracker API!"})
	})

	return router
}
