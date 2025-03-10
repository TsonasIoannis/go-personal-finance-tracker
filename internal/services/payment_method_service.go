package services

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

type PaymentMethodService struct {
	paymentMethodRepo PaymentMethodRepository
}

func NewPaymentMethodService(paymentMethodRepo PaymentMethodRepository) *PaymentMethodService {
	return &PaymentMethodService{paymentMethodRepo: paymentMethodRepo}
}

// AddPaymentMethod creates a new payment method
func (s *PaymentMethodService) AddPaymentMethod(paymentMethod *models.PaymentMethod) error {
	return s.paymentMethodRepo.CreatePaymentMethod(paymentMethod)
}

// GetPaymentMethodsByUser retrieves payment methods by user ID
func (s *PaymentMethodService) GetPaymentMethodsByUser(userID uint) ([]models.PaymentMethod, error) {
	return s.paymentMethodRepo.GetPaymentMethodsByUserID(userID)
}

// UpdatePaymentMethod modifies an existing payment method
func (s *PaymentMethodService) UpdatePaymentMethod(paymentMethod *models.PaymentMethod) error {
	return s.paymentMethodRepo.UpdatePaymentMethod(paymentMethod)
}

// DeletePaymentMethod removes a payment method
func (s *PaymentMethodService) DeletePaymentMethod(paymentMethodID uint) error {
	return s.paymentMethodRepo.DeletePaymentMethod(paymentMethodID)
}
