package routes

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userController *controllers.UserController, transactionController *controllers.TransactionController, budgetController *controllers.BudgetController) {
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)

	router.GET("/transactions", transactionController.GetTransactions)
	router.POST("/transactions", transactionController.CreateTransaction)
	router.DELETE("/transactions/:id", transactionController.DeleteTransaction)

	router.GET("/budgets", budgetController.GetBudgets)
	router.POST("/budgets", budgetController.CreateBudget)
	router.DELETE("/budgets/:id", budgetController.DeleteBudget)
}
