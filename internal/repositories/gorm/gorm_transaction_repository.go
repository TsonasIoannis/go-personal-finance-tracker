package repositories

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/filters"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
	"gorm.io/gorm"
)

// TransactionRepository defines the required repository methods
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransactionByID(ctx context.Context, id uint) (*models.Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID uint, filters filters.TransactionFilters) ([]models.Transaction, error)
	GetTransactionsPageByUserID(ctx context.Context, userID uint, params pagination.Params, filters filters.TransactionFilters) ([]models.Transaction, int64, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
	DeleteTransaction(ctx context.Context, id uint) error
}

// TransactionRepository handles DB operations for transactions
type GormTransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository initializes a new GormTransactionRepository
func NewTransactionRepository(db *gorm.DB) *GormTransactionRepository {
	return &GormTransactionRepository{db: db}
}

// CreateTransaction inserts a new transaction into the database
func (r *GormTransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// GetTransactionByID retrieves a transaction by its ID
func (r *GormTransactionRepository) GetTransactionByID(ctx context.Context, id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.WithContext(ctx).First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByUserID fetches all transactions for a specific user
func (r *GormTransactionRepository) GetTransactionsByUserID(ctx context.Context, userID uint, transactionFilters filters.TransactionFilters) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.transactionQuery(ctx, userID, transactionFilters).
		Order("date DESC").
		Order("id DESC").
		Find(&transactions).Error
	return transactions, err
}

// GetTransactionsPageByUserID fetches a paginated transaction list for a specific user.
func (r *GormTransactionRepository) GetTransactionsPageByUserID(ctx context.Context, userID uint, params pagination.Params, transactionFilters filters.TransactionFilters) ([]models.Transaction, int64, error) {
	var (
		transactions []models.Transaction
		total        int64
	)

	query := r.transactionQuery(ctx, userID, transactionFilters)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.transactionQuery(ctx, userID, transactionFilters).
		Order("date DESC").
		Order("id DESC").
		Offset(params.Offset()).
		Limit(params.PageSize).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// UpdateTransaction updates an existing transaction
func (r *GormTransactionRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

// DeleteTransaction removes a transaction from the database
func (r *GormTransactionRepository) DeleteTransaction(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Transaction{}, id).Error
}

func (r *GormTransactionRepository) transactionQuery(ctx context.Context, userID uint, transactionFilters filters.TransactionFilters) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&models.Transaction{}).Where("user_id = ?", userID)

	if transactionFilters.Type != "" {
		query = query.Where("\"type\" = ?", transactionFilters.Type)
	}

	if transactionFilters.CategoryID != nil {
		query = query.Where("category_id = ?", *transactionFilters.CategoryID)
	}

	if transactionFilters.From != nil {
		query = query.Where("date >= ?", *transactionFilters.From)
	}

	if transactionFilters.To != nil {
		query = query.Where("date <= ?", *transactionFilters.To)
	}

	return query
}
