package entity

type OrderItemEntity struct {
	ID           int64
	OrderID      int64
	ProductID    int64
	Quantity     int64
	OrderCode    string
	ProductName  string
	ProductImage string
	Price        int64
}
