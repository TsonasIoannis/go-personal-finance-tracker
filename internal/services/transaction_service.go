package services

import (
	"errors"

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

type TransactionService struct {
	transactionRepo TransactionRepository
	budgetRepo      BudgetRepository
}

func NewTransactionService(transactionRepo TransactionRepository, budgetRepo BudgetRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo, budgetRepo: budgetRepo}
}

// AddTransaction validates and saves a transaction
func (s *TransactionService) AddTransaction(transaction *models.Transaction) error {
	// Check if transaction exceeds budget
	budgets, err := s.budgetRepo.GetBudgetsByUserID(transaction.UserID)
	if err != nil {
		return err
	}

	for _, budget := range budgets {
		if budget.CategoryID == transaction.CategoryID && transaction.Type == "expense" {
			if transaction.Amount > budget.Limit {
				return errors.New("transaction exceeds budget limit")
			}
		}
	}

	return s.transactionRepo.CreateTransaction(transaction)
}

// GetTransactionsByUser retrieves all transactions for a user
func (s *TransactionService) GetTransactionsByUser(userID uint) ([]models.Transaction, error) {
	return s.transactionRepo.GetTransactionsByUserID(userID)
}

// DeleteTransaction removes a transaction by ID
func (s *TransactionService) DeleteTransaction(transactionID uint) error {
	return s.transactionRepo.DeleteTransaction(transactionID)
}
