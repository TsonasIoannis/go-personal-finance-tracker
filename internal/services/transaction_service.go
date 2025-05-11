package services

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// TransactionService defines the interface for transaction operations
type TransactionService interface {
	AddTransaction(transaction *models.Transaction) error
	GetTransactionsByUser(userID uint) ([]models.Transaction, error)
	DeleteTransaction(transactionID uint) error
}
