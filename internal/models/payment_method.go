package models

type PaymentMethod struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"size:50;uniqueIndex;not null"` // "Credit Card", "Cash", "PayPal"
	UserID uint   `gorm:"not null;index"`
}
