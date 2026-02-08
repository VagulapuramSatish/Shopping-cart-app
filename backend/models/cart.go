package models

import "time"

type Cart struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"uniqueIndex"`
	Name      string
	Status    string    `gorm:"default:active"`
	CreatedAt time.Time
}
