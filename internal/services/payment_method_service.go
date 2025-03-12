package services

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// PaymentMethodService defines the interface for payment method operations
type PaymentMethodService interface {
	AddPaymentMethod(paymentMethod *models.PaymentMethod) error
	GetPaymentMethodsByUser(userID uint) ([]models.PaymentMethod, error)
	UpdatePaymentMethod(paymentMethod *models.PaymentMethod) error
	DeletePaymentMethod(paymentMethodID uint) error
}
