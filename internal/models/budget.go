package models

import "time"

type Budget struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null;index"`
	CategoryID uint      `gorm:"not null;index"`
	Limit      float64   `gorm:"not null"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
}
