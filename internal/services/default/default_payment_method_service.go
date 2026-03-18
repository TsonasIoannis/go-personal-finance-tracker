package services

import (
	"context"

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
func (s *DefaultPaymentMethodService) AddPaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	return s.paymentMethodRepo.CreatePaymentMethod(ctx, paymentMethod)
}

// GetPaymentMethodsByUser retrieves payment methods by user ID
func (s *DefaultPaymentMethodService) GetPaymentMethodsByUser(ctx context.Context, userID uint) ([]models.PaymentMethod, error) {
	return s.paymentMethodRepo.GetPaymentMethodsByUserID(ctx, userID)
}

// UpdatePaymentMethod modifies an existing payment method
func (s *DefaultPaymentMethodService) UpdatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	return s.paymentMethodRepo.UpdatePaymentMethod(ctx, paymentMethod)
}

// DeletePaymentMethod removes a payment method
func (s *DefaultPaymentMethodService) DeletePaymentMethod(ctx context.Context, paymentMethodID uint) error {
	return s.paymentMethodRepo.DeletePaymentMethod(ctx, paymentMethodID)
}
