package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
)

// TransactionRepository handles DB operations for transactions
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository initializes a new TransactionRepository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction inserts a new transaction into the database
func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// GetTransactionByID retrieves a transaction by its ID
func (r *TransactionRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByUserID fetches all transactions for a specific user
func (r *TransactionRepository) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

// UpdateTransaction updates an existing transaction
func (r *TransactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

// DeleteTransaction removes a transaction from the database
func (r *TransactionRepository) DeleteTransaction(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}
