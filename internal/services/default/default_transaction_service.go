package services

import (
	"errors"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories"
)

type DefaultTransactionService struct {
	transactionRepo repositories.TransactionRepository
	budgetRepo      repositories.BudgetRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, budgetRepo repositories.BudgetRepository) *DefaultTransactionService {
	return &DefaultTransactionService{transactionRepo: transactionRepo, budgetRepo: budgetRepo}
}

// AddTransaction validates and saves a transaction
func (s *DefaultTransactionService) AddTransaction(transaction *models.Transaction) error {
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
func (s *DefaultTransactionService) GetTransactionsByUser(userID uint) ([]models.Transaction, error) {
	return s.transactionRepo.GetTransactionsByUserID(userID)
}

// DeleteTransaction removes a transaction by ID
func (s *DefaultTransactionService) DeleteTransaction(transactionID uint) error {
	return s.transactionRepo.DeleteTransaction(transactionID)
}
