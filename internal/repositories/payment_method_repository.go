package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// PaymentMethodRepository defines the required repository methods
type PaymentMethodRepository interface {
	CreatePaymentMethod(paymentMethod *models.PaymentMethod) error
	GetPaymentMethodByID(id uint) (*models.PaymentMethod, error)
	GetPaymentMethodsByUserID(userID uint) ([]models.PaymentMethod, error)
	UpdatePaymentMethod(paymentMethod *models.PaymentMethod) error
	DeletePaymentMethod(id uint) error
}
