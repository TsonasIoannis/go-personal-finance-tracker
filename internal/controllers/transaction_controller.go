package controllers

import (
	"net/http"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}

// CreateTransaction adds a new transaction
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := tc.transactionService.AddTransaction(&transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction added"})
}

// GetTransactions fetches all transactions for a user
func (tc *TransactionController) GetTransactions(c *gin.Context) {
	userID := uint(1) // Later, extract from JWT or request context

	transactions, err := tc.transactionService.GetTransactionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// DeleteTransaction removes a transaction
func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	transactionID := uint(1) // Extract from URL params in real implementation

	err := tc.transactionService.DeleteTransaction(transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}
