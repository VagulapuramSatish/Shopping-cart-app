package models

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	CartID    uint
	UserID    uint
	CreatedAt time.Time
}
