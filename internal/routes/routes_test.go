package routes

import (
	"context"
	"sort"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/auth"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/controllers"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
	"github.com/gin-gonic/gin"
)

type stubUserService struct{}

func (stubUserService) RegisterUser(context.Context, string, string, string) (*models.User, error) {
	return nil, nil
}

func (stubUserService) AuthenticateUser(context.Context, string, string) (*models.User, error) {
	return nil, nil
}

type stubTransactionService struct{}

func (stubTransactionService) AddTransaction(context.Context, *models.Transaction) error {
	return nil
}

func (stubTransactionService) GetTransactionsByUser(context.Context, uint) ([]models.Transaction, error) {
	return nil, nil
}

func (stubTransactionService) GetTransactionsPageByUser(context.Context, uint, pagination.Params) ([]models.Transaction, int64, error) {
	return nil, 0, nil
}

func (stubTransactionService) DeleteTransactionForUser(context.Context, uint, uint) error {
	return nil
}

type stubBudgetService struct{}

func (stubBudgetService) CreateBudget(context.Context, *models.Budget) error {
	return nil
}

func (stubBudgetService) UpdateBudget(context.Context, *models.Budget) error {
	return nil
}

func (stubBudgetService) GetBudgetsByUser(context.Context, uint) ([]models.Budget, error) {
	return nil, nil
}

func (stubBudgetService) GetBudgetsPageByUser(context.Context, uint, pagination.Params) ([]models.Budget, int64, error) {
	return nil, 0, nil
}

func (stubBudgetService) DeleteBudgetForUser(context.Context, uint, uint) error {
	return nil
}

type stubTokenManager struct{}

func (stubTokenManager) GenerateToken(*models.User) (string, error) {
	return "", nil
}

func (stubTokenManager) ParseToken(string) (*auth.Claims, error) {
	return &auth.Claims{UserID: 1, Email: "test@example.com"}, nil
}

func TestSetupRoutesRegistersLegacyAndVersionedAPIPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	userController := controllers.NewUserController(stubUserService{}, stubTokenManager{})
	transactionController := controllers.NewTransactionController(stubTransactionService{})
	budgetController := controllers.NewBudgetController(stubBudgetService{})

	SetupRoutes(router, func(c *gin.Context) { c.Next() }, userController, transactionController, budgetController)

	got := make([]string, 0, len(router.Routes()))
	for _, route := range router.Routes() {
		got = append(got, route.Method+" "+route.Path)
	}
	sort.Strings(got)

	want := []string{
		"DELETE /api/v1/budgets/:id",
		"DELETE /api/v1/transactions/:id",
		"DELETE /budgets/:id",
		"DELETE /transactions/:id",
		"GET /api/v1/budgets",
		"GET /api/v1/transactions",
		"GET /budgets",
		"GET /transactions",
		"POST /api/v1/budgets",
		"POST /api/v1/login",
		"POST /api/v1/register",
		"POST /api/v1/transactions",
		"POST /budgets",
		"POST /login",
		"POST /register",
		"POST /transactions",
	}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("expected %d routes, got %d: %v", len(want), len(got), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected route %q at index %d, got %q", want[i], i, got[i])
		}
	}
}
