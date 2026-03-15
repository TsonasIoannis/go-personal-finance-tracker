package controllers

import (
	"net/http"
	"strconv"
	"time"

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

type createTransactionRequest struct {
	Type       string    `json:"type" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	CategoryID uint      `json:"category_id" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
	Note       string    `json:"note"`
}

// CreateTransaction adds a new transaction
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	transaction := models.Transaction{
		UserID:     userID,
		Type:       req.Type,
		Amount:     req.Amount,
		CategoryID: req.CategoryID,
		Date:       req.Date,
		Note:       req.Note,
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
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	transactions, err := tc.transactionService.GetTransactionsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// DeleteTransaction removes a transaction
func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction id"})
		return
	}

	err = tc.transactionService.DeleteTransactionForUser(userID, uint(transactionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}
