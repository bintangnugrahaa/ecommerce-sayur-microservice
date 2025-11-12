package entity

import "time"

type OrderEntity struct {
	ID            int64
	OrderCode     string
	BuyerID       int64
	OrderDate     string
	Status        string
	TotalAmount   int64
	PaymentMethod string
	ShippingType  string
	ShippingFee   int64
	OrderTime     string
	Remarks       string
	CreatedAt     time.Time
	OrderItems    []OrderItemEntity
	BuyerName     string
	BuyerEmail    string
	BuyerPhone    string
	BuyerAddress  string
	BuyerLat      string
	BuyerLng      string
}

type QueryStringEntity struct {
	Page    int64
	Search  string
	Limit   int64
	Status  string
	BuyerID int64
}
