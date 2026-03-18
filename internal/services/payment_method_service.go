package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// PaymentMethodService defines the interface for payment method operations
type PaymentMethodService interface {
	AddPaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error
	GetPaymentMethodsByUser(ctx context.Context, userID uint) ([]models.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, paymentMethodID uint) error
}
