package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/filters"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
)

// TransactionService defines the interface for transaction operations
type TransactionService interface {
	AddTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactionsByUser(ctx context.Context, userID uint, filters filters.TransactionFilters) ([]models.Transaction, error)
	GetTransactionsPageByUser(ctx context.Context, userID uint, params pagination.Params, filters filters.TransactionFilters) ([]models.Transaction, int64, error)
	DeleteTransactionForUser(ctx context.Context, userID, transactionID uint) error
}
