package services

import (
	"context"
	"errors"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPaymentMethodRepository implements the PaymentMethodRepository interface
type MockPaymentMethodRepository struct {
	mock.Mock
}

func (m *MockPaymentMethodRepository) CreatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	args := m.Called(ctx, paymentMethod)
	return args.Error(0)
}

func (m *MockPaymentMethodRepository) GetPaymentMethodByID(ctx context.Context, id uint) (*models.PaymentMethod, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.PaymentMethod), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPaymentMethodRepository) GetPaymentMethodsByUserID(ctx context.Context, userID uint) ([]models.PaymentMethod, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodRepository) UpdatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	args := m.Called(ctx, paymentMethod)
	return args.Error(0)
}

func (m *MockPaymentMethodRepository) DeletePaymentMethod(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestAddPaymentMethod(t *testing.T) {
	mockRepo := new(MockPaymentMethodRepository)
	service := NewPaymentMethodService(mockRepo)
	ctx := context.Background()

	t.Run("Create valid payment method", func(t *testing.T) {
		paymentMethod := &models.PaymentMethod{Name: "Credit Card", UserID: 1}
		mockRepo.On("CreatePaymentMethod", ctx, paymentMethod).Return(nil)

		err := service.AddPaymentMethod(ctx, paymentMethod)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPaymentMethodsByUser(t *testing.T) {
	mockRepo := new(MockPaymentMethodRepository)
	service := NewPaymentMethodService(mockRepo)
	ctx := context.Background()

	t.Run("Retrieve multiple payment methods", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset expectations

		paymentMethods := []models.PaymentMethod{
			{ID: 1, Name: "Apple Pay", UserID: 1},
			{ID: 2, Name: "Google Pay", UserID: 1},
		}

		mockRepo.On("GetPaymentMethodsByUserID", ctx, uint(1)).Return(paymentMethods, nil)

		result, err := service.GetPaymentMethodsByUser(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Retrieve payment methods when none exist", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset expectations

		mockRepo.On("GetPaymentMethodsByUserID", ctx, uint(999)).Return([]models.PaymentMethod{}, nil)

		result, err := service.GetPaymentMethodsByUser(ctx, 999)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdatePaymentMethod(t *testing.T) {
	mockRepo := new(MockPaymentMethodRepository)
	service := NewPaymentMethodService(mockRepo)
	ctx := context.Background()

	t.Run("Update existing payment method", func(t *testing.T) {
		paymentMethod := &models.PaymentMethod{ID: 1, Name: "Bank Transfer", UserID: 1}
		mockRepo.On("UpdatePaymentMethod", ctx, paymentMethod).Return(nil)

		err := service.UpdatePaymentMethod(ctx, paymentMethod)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to update non-existent payment method", func(t *testing.T) {
		paymentMethod := &models.PaymentMethod{ID: 9999, Name: "Cryptocurrency", UserID: 1}
		mockRepo.On("UpdatePaymentMethod", ctx, paymentMethod).Return(errors.New("payment method not found"))

		err := service.UpdatePaymentMethod(ctx, paymentMethod)
		assert.Error(t, err)
		assert.Equal(t, "payment method not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestDeletePaymentMethod(t *testing.T) {
	mockRepo := new(MockPaymentMethodRepository)
	service := NewPaymentMethodService(mockRepo)
	ctx := context.Background()

	t.Run("Delete existing payment method", func(t *testing.T) {
		mockRepo.On("DeletePaymentMethod", ctx, uint(1)).Return(nil)

		err := service.DeletePaymentMethod(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete non-existent payment method", func(t *testing.T) {
		mockRepo.On("DeletePaymentMethod", ctx, uint(9999)).Return(errors.New("payment method not found"))

		err := service.DeletePaymentMethod(ctx, 9999)
		assert.Error(t, err)
		assert.Equal(t, "payment method not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
