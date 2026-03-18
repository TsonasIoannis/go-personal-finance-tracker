package app

import (
	"log/slog"
	"net/http"

	_ "github.com/TsonasIoannis/go-personal-finance-tracker/docs"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/auth"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/config"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/controllers"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/handlers"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/middleware"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/observability"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/persistence"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/routes"
	services "github.com/TsonasIoannis/go-personal-finance-tracker/internal/services/default"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(cfg config.Config, db database.Database, repositories persistence.Repositories) *http.Server {
	router := newRouter(cfg, db, repositories)

	return &http.Server{
		Addr:              cfg.Address(),
		Handler:           router,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
	}
}

func newRouter(cfg config.Config, db database.Database, repositories persistence.Repositories) *gin.Engine {
	metrics := observability.NewHTTPMetrics()

	router := gin.New()
	router.Use(
		middleware.RequestIDMiddleware(),
		observability.TracingMiddleware(),
		middleware.StructuredLoggerMiddleware(slog.Default()),
		metrics.Middleware(),
		middleware.RecoveryMiddleware(slog.Default()),
	)

	tokenManager := auth.NewJWTManager(cfg.JWTSecret, cfg.Auth.TokenTTL)
	authMiddleware := middleware.AuthMiddleware(tokenManager)

	userService := services.NewUserService(repositories.Users)
	transactionService := services.NewTransactionService(repositories.Transactions, repositories.Budgets)
	budgetService := services.NewBudgetService(repositories.Budgets)

	userController := controllers.NewUserController(userService, tokenManager)
	transactionController := controllers.NewTransactionController(transactionService)
	budgetController := controllers.NewBudgetController(budgetService)

	routes.SetupRoutes(router, authMiddleware, userController, transactionController, budgetController)

	router.GET("/metrics", gin.WrapH(metrics.Handler()))
	router.GET("/health", handlers.HealthCheckHandler)
	router.GET("/ready", handlers.ReadinessCheckHandler(db))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Personal Finance Tracker API!"})
	})

	return router
}
