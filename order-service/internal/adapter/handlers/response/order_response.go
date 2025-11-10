package response

type OrderAdminList struct {
	ID            int64  `json:"id"`
	OrderCode     string `json:"order_code"`
	ProductImage  string `json:"product_image"`
	CustomerName  string `json:"customer_name"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	TotalAmount   int64  `json:"total_amount"`
}
