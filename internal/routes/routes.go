package routes

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc, userController *controllers.UserController, transactionController *controllers.TransactionController, budgetController *controllers.BudgetController) {
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)

	protected := router.Group("/")
	protected.Use(authMiddleware)
	protected.GET("/transactions", transactionController.GetTransactions)
	protected.POST("/transactions", transactionController.CreateTransaction)
	protected.DELETE("/transactions/:id", transactionController.DeleteTransaction)
	protected.GET("/budgets", budgetController.GetBudgets)
	protected.POST("/budgets", budgetController.CreateBudget)
	protected.DELETE("/budgets/:id", budgetController.DeleteBudget)
}
