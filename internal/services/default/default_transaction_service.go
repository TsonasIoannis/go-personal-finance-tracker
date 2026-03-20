package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/filters"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
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
func (s *DefaultTransactionService) AddTransaction(ctx context.Context, transaction *models.Transaction) error {
	// Check if transaction exceeds budget
	budgets, err := s.budgetRepo.GetBudgetsByUserID(ctx, transaction.UserID)
	if err != nil {
		return apperrors.Internal("budget_lookup_failed", "failed to validate transaction budget", err)
	}

	for _, budget := range budgets {
		if budget.CategoryID == transaction.CategoryID && transaction.Type == "expense" {
			if transaction.Amount > budget.Limit {
				return apperrors.Validation("budget_limit_exceeded", "transaction exceeds budget limit")
			}
		}
	}

	if err := s.transactionRepo.CreateTransaction(ctx, transaction); err != nil {
		return apperrors.Internal("transaction_create_failed", "failed to create transaction", err)
	}

	return nil
}

// GetTransactionsByUser retrieves all transactions for a user
func (s *DefaultTransactionService) GetTransactionsByUser(ctx context.Context, userID uint, transactionFilters filters.TransactionFilters) ([]models.Transaction, error) {
	transactions, err := s.transactionRepo.GetTransactionsByUserID(ctx, userID, transactionFilters)
	if err != nil {
		return nil, apperrors.Internal("transactions_fetch_failed", "failed to retrieve transactions", err)
	}

	return transactions, nil
}

// GetTransactionsPageByUser retrieves a paginated transaction list for a user.
func (s *DefaultTransactionService) GetTransactionsPageByUser(ctx context.Context, userID uint, params pagination.Params, transactionFilters filters.TransactionFilters) ([]models.Transaction, int64, error) {
	transactions, total, err := s.transactionRepo.GetTransactionsPageByUserID(ctx, userID, params, transactionFilters)
	if err != nil {
		return nil, 0, apperrors.Internal("transactions_fetch_failed", "failed to retrieve transactions", err)
	}

	return transactions, total, nil
}

// DeleteTransactionForUser removes a transaction that belongs to the authenticated user.
func (s *DefaultTransactionService) DeleteTransactionForUser(ctx context.Context, userID, transactionID uint) error {
	transaction, err := s.transactionRepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return apperrors.NotFound("transaction_not_found", "transaction not found")
	}

	if transaction.UserID != userID {
		return apperrors.NotFound("transaction_not_found", "transaction not found")
	}

	if err := s.transactionRepo.DeleteTransaction(ctx, transactionID); err != nil {
		return apperrors.Internal("transaction_delete_failed", "failed to delete transaction", err)
	}

	return nil
}
