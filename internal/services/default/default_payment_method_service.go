package services

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories"
)

type DefaultPaymentMethodService struct {
	paymentMethodRepo repositories.PaymentMethodRepository
}

func NewPaymentMethodService(paymentMethodRepo repositories.PaymentMethodRepository) *DefaultPaymentMethodService {
	return &DefaultPaymentMethodService{paymentMethodRepo: paymentMethodRepo}
}

// AddPaymentMethod creates a new payment method
func (s *DefaultPaymentMethodService) AddPaymentMethod(paymentMethod *models.PaymentMethod) error {
	return s.paymentMethodRepo.CreatePaymentMethod(paymentMethod)
}

// GetPaymentMethodsByUser retrieves payment methods by user ID
func (s *DefaultPaymentMethodService) GetPaymentMethodsByUser(userID uint) ([]models.PaymentMethod, error) {
	return s.paymentMethodRepo.GetPaymentMethodsByUserID(userID)
}

// UpdatePaymentMethod modifies an existing payment method
func (s *DefaultPaymentMethodService) UpdatePaymentMethod(paymentMethod *models.PaymentMethod) error {
	return s.paymentMethodRepo.UpdatePaymentMethod(paymentMethod)
}

// DeletePaymentMethod removes a payment method
func (s *DefaultPaymentMethodService) DeletePaymentMethod(paymentMethodID uint) error {
	return s.paymentMethodRepo.DeletePaymentMethod(paymentMethodID)
}
