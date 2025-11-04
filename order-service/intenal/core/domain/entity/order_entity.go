package entity

import "time"

type OrderEhtity struct {
	ID           int64
	OrderCode    string
	BuyerID      int64
	OrderDate    string
	Status       string
	TotalAmount  int64
	ShippingType string
	ShippingFee  int64
	OrderTime    string
	Remarks      string
	CreatedAt    time.Time
	OrderItems   []OrderItemEntity
	BuyerName    string
}
