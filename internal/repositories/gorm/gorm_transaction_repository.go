package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
)

// TransactionRepository defines the required repository methods
type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactionsByUserID(userID uint) ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uint) error
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
func (r *GormTransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

// GetTransactionByID retrieves a transaction by its ID
func (r *GormTransactionRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByUserID fetches all transactions for a specific user
func (r *GormTransactionRepository) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

// UpdateTransaction updates an existing transaction
func (r *GormTransactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

// DeleteTransaction removes a transaction from the database
func (r *GormTransactionRepository) DeleteTransaction(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}
