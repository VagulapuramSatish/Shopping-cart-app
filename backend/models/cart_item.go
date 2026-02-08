package models

type CartItem struct {
	CartID uint `gorm:"primaryKey"`
	ItemID uint `gorm:"primaryKey"`
}
