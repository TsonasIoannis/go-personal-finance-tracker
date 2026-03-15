package models

import (
	"time"
)

type User struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"size:100;not null" json:"name"`
	Email        string        `gorm:"uniqueIndex;not null" json:"email"`
	Password     string        `gorm:"not null" json:"-"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}
