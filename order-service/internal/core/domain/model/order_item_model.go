package model

import "time"

type OrderItem struct {
	ID        int64 `gorm:"primaryKey"`
	OrderID   int64 `gorm:"order_id"`
	ProductID int64 `gorm:"product_id"`
	Quantity  int64 `gorm:"quantity"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Order     Order `gorm:"foreignKey:OrderID;references:ID"`
}
