package persistence

import (
	repositorycontracts "github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories"
	gormrepositories "github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories/gorm"
	"gorm.io/gorm"
)

type Repositories struct {
	Users        repositorycontracts.UserRepository
	Transactions repositorycontracts.TransactionRepository
	Budgets      repositorycontracts.BudgetRepository
}

func NewGormRepositories(db *gorm.DB) Repositories {
	return Repositories{
		Users:        gormrepositories.NewUserRepository(db),
		Transactions: gormrepositories.NewTransactionRepository(db),
		Budgets:      gormrepositories.NewGormBudgetRepository(db),
	}
}
