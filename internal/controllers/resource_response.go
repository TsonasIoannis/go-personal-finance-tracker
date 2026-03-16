package controllers

import (
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

type transactionResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Type       string    `json:"type"`
	Amount     float64   `json:"amount"`
	CategoryID uint      `json:"category_id"`
	Date       time.Time `json:"date"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type budgetResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	CategoryID uint      `json:"category_id"`
	Limit      float64   `json:"limit"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

func newTransactionResponse(transaction models.Transaction) transactionResponse {
	return transactionResponse{
		ID:         transaction.ID,
		UserID:     transaction.UserID,
		Type:       transaction.Type,
		Amount:     transaction.Amount,
		CategoryID: transaction.CategoryID,
		Date:       transaction.Date,
		Note:       transaction.Note,
		CreatedAt:  transaction.CreatedAt,
		UpdatedAt:  transaction.UpdatedAt,
	}
}

func newBudgetResponse(budget models.Budget) budgetResponse {
	return budgetResponse{
		ID:         budget.ID,
		UserID:     budget.UserID,
		CategoryID: budget.CategoryID,
		Limit:      budget.Limit,
		StartDate:  budget.StartDate,
		EndDate:    budget.EndDate,
	}
}

func newTransactionResponses(transactions []models.Transaction) []transactionResponse {
	responses := make([]transactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		responses = append(responses, newTransactionResponse(transaction))
	}
	return responses
}

func newBudgetResponses(budgets []models.Budget) []budgetResponse {
	responses := make([]budgetResponse, 0, len(budgets))
	for _, budget := range budgets {
		responses = append(responses, newBudgetResponse(budget))
	}
	return responses
}
