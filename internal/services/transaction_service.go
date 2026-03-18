package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// TransactionService defines the interface for transaction operations
type TransactionService interface {
	AddTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactionsByUser(ctx context.Context, userID uint) ([]models.Transaction, error)
	DeleteTransactionForUser(ctx context.Context, userID, transactionID uint) error
}
