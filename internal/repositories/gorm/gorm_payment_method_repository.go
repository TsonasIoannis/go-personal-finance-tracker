package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
)

// PaymentMethodRepository defines the required repository methods
type PaymentMethodRepository interface {
	CreatePaymentMethod(paymentMethod *models.PaymentMethod) error
	GetPaymentMethodByID(id uint) (*models.PaymentMethod, error)
	GetPaymentMethodsByUserID(userID uint) ([]models.PaymentMethod, error)
	UpdatePaymentMethod(paymentMethod *models.PaymentMethod) error
	DeletePaymentMethod(id uint) error
}

// PaymentMethodRepository handles DB operations for payment methods
type GormPaymentMethodRepository struct {
	db *gorm.DB
}

// NewPaymentMethodRepository initializes a new GormPaymentMethodRepository
func NewPaymentMethodRepository(db *gorm.DB) *GormPaymentMethodRepository {
	return &GormPaymentMethodRepository{db: db}
}

// CreatePaymentMethod inserts a new payment method into the database
func (r *GormPaymentMethodRepository) CreatePaymentMethod(paymentMethod *models.PaymentMethod) error {
	return r.db.Create(paymentMethod).Error
}

// GetPaymentMethodByID retrieves a payment method by its ID
func (r *GormPaymentMethodRepository) GetPaymentMethodByID(id uint) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	err := r.db.First(&paymentMethod, id).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

// GetPaymentMethodsByUserID fetches all payment methods for a specific user
func (r *GormPaymentMethodRepository) GetPaymentMethodsByUserID(userID uint) ([]models.PaymentMethod, error) {
	var paymentMethods []models.PaymentMethod
	err := r.db.Where("user_id = ?", userID).Find(&paymentMethods).Error
	return paymentMethods, err
}

// UpdatePaymentMethod updates an existing payment method
func (r *GormPaymentMethodRepository) UpdatePaymentMethod(paymentMethod *models.PaymentMethod) error {
	return r.db.Save(paymentMethod).Error
}

// DeletePaymentMethod removes a payment method from the database
func (r *GormPaymentMethodRepository) DeletePaymentMethod(id uint) error {
	return r.db.Delete(&models.PaymentMethod{}, id).Error
}
