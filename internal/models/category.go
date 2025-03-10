package models

type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:100;uniqueIndex;not null"`
	Description string `gorm:"size:255"`
}
