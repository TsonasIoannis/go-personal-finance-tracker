package models

import (
	"time"
)

type Transaction struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null;index"`
	Type       string    `gorm:"size:10;not null"` // "income" or "expense"
	Amount     float64   `gorm:"not null"`
	CategoryID uint      `gorm:"not null;index"`
	Date       time.Time `gorm:"not null"`
	Note       string    `gorm:"size:255"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type RecurringTransaction struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"not null;index"`
	CategoryID  uint       `gorm:"not null;index"`
	Amount      float64    `gorm:"not null"`
	Frequency   string     `gorm:"size:20;not null"` // "daily", "weekly", "monthly"
	NextDueDate time.Time  `gorm:"not null"`
	EndDate     *time.Time // Nullable - stops recurrence
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
