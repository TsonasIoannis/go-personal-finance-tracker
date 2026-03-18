package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/httpapi"
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
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.WriteError(c, apperrors.Validation("invalid_request", "invalid request payload"))
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

	err := tc.transactionService.AddTransaction(ctx, &transaction)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction added"})
}

// GetTransactions fetches all transactions for a user
func (tc *TransactionController) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	transactions, err := tc.transactionService.GetTransactionsByUser(ctx, userID)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, newTransactionResponses(transactions))
}

// GetTransactionsPage fetches a paginated transaction list for a user.
func (tc *TransactionController) GetTransactionsPage(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	params, err := parsePaginationParams(c)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	transactions, total, err := tc.transactionService.GetTransactionsPageByUser(ctx, userID, params)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, paginatedResponse[transactionResponse]{
		Data:       newTransactionResponses(transactions),
		Pagination: newPaginationResponse(params, total),
	})
}

// DeleteTransaction removes a transaction
func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.WriteError(c, apperrors.Validation("invalid_transaction_id", "invalid transaction id"))
		return
	}

	err = tc.transactionService.DeleteTransactionForUser(ctx, userID, uint(transactionID))
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}
