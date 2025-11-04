package model

import "time"

type Order struct {
	ID           int64     `gorm:"primaryKey"`
	OrderCode    string    `gorm:"oder_code"`
	BuyerID      int64     `gorm:"buyer_id"`
	OrderDate    time.Time `gorm:"order_date"`
	Status       string    `gorm:"status"`
	TotalAmount  int64     `gorm:"total_amount"`
	ShippingType string    `gorm:"shipping_type"`
	ShippingFee  int64     `gorm:"shipping_fee"`
	OrderTime    time.Time `gorm:"order_time"`
	Remarks      string    `gorm:"remarks"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}
