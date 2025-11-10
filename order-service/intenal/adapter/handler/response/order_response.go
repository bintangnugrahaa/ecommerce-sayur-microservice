package response

type OrderAdminList struct {
	ID            string  `json:"id"`
	OrderCode     string  `json:"order_code"`
	ProductImage  string  `json:"product_image"`
	CustomeName   string  `json:"customer_name"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	TotalAmount   float64 `json:"total_amount"`
}
