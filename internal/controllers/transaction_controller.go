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
// @Summary Create a transaction
// @Description Create a transaction for the authenticated user.
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body createTransactionRequest true "Transaction payload"
// @Success 201 {object} messageResponse
// @Failure 400 {object} httpapi.ErrorResponse
// @Failure 401 {object} httpapi.ErrorResponse
// @Failure 500 {object} httpapi.ErrorResponse
// @Router /api/v1/transactions [post]
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

	transactionFilters, err := parseTransactionFilters(c)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	transactions, err := tc.transactionService.GetTransactionsByUser(ctx, userID, transactionFilters)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, newTransactionResponses(transactions))
}

// GetTransactionsPage fetches a paginated transaction list for a user.
// @Summary List transactions
// @Description List the authenticated user's transactions with pagination and optional filtering.
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" minimum(1)
// @Param page_size query int false "Items per page" minimum(1) maximum(100)
// @Param type query string false "Transaction type" Enums(income, expense)
// @Param category_id query int false "Category ID" minimum(1)
// @Param from query string false "Start date/time filter (RFC3339 or YYYY-MM-DD)"
// @Param to query string false "End date/time filter (RFC3339 or YYYY-MM-DD)"
// @Success 200 {object} transactionPageResponse
// @Failure 400 {object} httpapi.ErrorResponse
// @Failure 401 {object} httpapi.ErrorResponse
// @Failure 500 {object} httpapi.ErrorResponse
// @Router /api/v1/transactions [get]
func (tc *TransactionController) GetTransactionsPage(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	transactionFilters, err := parseTransactionFilters(c)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	params, err := parsePaginationParams(c)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	transactions, total, err := tc.transactionService.GetTransactionsPageByUser(ctx, userID, params, transactionFilters)
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
// @Summary Delete a transaction
// @Description Delete one of the authenticated user's transactions.
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID" minimum(1)
// @Success 200 {object} messageResponse
// @Failure 400 {object} httpapi.ErrorResponse
// @Failure 401 {object} httpapi.ErrorResponse
// @Failure 404 {object} httpapi.ErrorResponse
// @Failure 500 {object} httpapi.ErrorResponse
// @Router /api/v1/transactions/{id} [delete]
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
