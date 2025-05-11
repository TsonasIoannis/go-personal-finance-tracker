package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// TransactionRepository defines the required repository methods
type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactionsByUserID(userID uint) ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uint) error
}
