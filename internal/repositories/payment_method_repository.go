package repositories

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// PaymentMethodRepository defines the required repository methods
type PaymentMethodRepository interface {
	CreatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error
	GetPaymentMethodByID(ctx context.Context, id uint) (*models.PaymentMethod, error)
	GetPaymentMethodsByUserID(ctx context.Context, userID uint) ([]models.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, id uint) error
}
