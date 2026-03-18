package routes

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc, userController *controllers.UserController, transactionController *controllers.TransactionController, budgetController *controllers.BudgetController) {
	registerPublicRoutes(router, userController)
	legacyProtected := router.Group("/")
	legacyProtected.Use(authMiddleware)
	registerLegacyProtectedRoutes(legacyProtected, transactionController, budgetController)

	v1 := router.Group("/api/v1")
	registerPublicRoutes(v1, userController)
	v1Protected := v1.Group("/")
	v1Protected.Use(authMiddleware)
	registerVersionedProtectedRoutes(v1Protected, transactionController, budgetController)
}

func registerPublicRoutes(router gin.IRoutes, userController *controllers.UserController) {
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
}

func registerLegacyProtectedRoutes(router gin.IRoutes, transactionController *controllers.TransactionController, budgetController *controllers.BudgetController) {
	router.GET("/transactions", transactionController.GetTransactions)
	router.POST("/transactions", transactionController.CreateTransaction)
	router.DELETE("/transactions/:id", transactionController.DeleteTransaction)
	router.GET("/budgets", budgetController.GetBudgets)
	router.POST("/budgets", budgetController.CreateBudget)
	router.DELETE("/budgets/:id", budgetController.DeleteBudget)
}

func registerVersionedProtectedRoutes(router gin.IRoutes, transactionController *controllers.TransactionController, budgetController *controllers.BudgetController) {
	router.GET("/transactions", transactionController.GetTransactionsPage)
	router.POST("/transactions", transactionController.CreateTransaction)
	router.DELETE("/transactions/:id", transactionController.DeleteTransaction)
	router.GET("/budgets", budgetController.GetBudgetsPage)
	router.POST("/budgets", budgetController.CreateBudget)
	router.DELETE("/budgets/:id", budgetController.DeleteBudget)
}
