package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Token     string    `gorm:"type:text"`
	CartID    *uint
	CreatedAt time.Time
}

