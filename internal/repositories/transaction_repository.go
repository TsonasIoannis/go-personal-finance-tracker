package repositories

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// TransactionRepository defines the required repository methods
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactionByID(ctx context.Context, id uint) (*models.Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID uint) ([]models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
	DeleteTransaction(ctx context.Context, id uint) error
}
