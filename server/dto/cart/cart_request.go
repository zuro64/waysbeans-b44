package cartdto

type CreateCartRequest struct {
	TransactionID int  `json:"transaction_id" form:"transaction_id" validate:"required"`
	ProductID     int  `json:"product_id" form:"product_id" validate:"required"`
	OrderQuantity int  `json:"order_quantity" form:"order_quantity" validate:"required"`
	Checkout      bool `json:"checkout" form:"checkout"`
}

type UpdateCartRequest struct {
	OrderQuantity int  `json:"order_quantity" form:"order_quantity"`
	Checkout      bool `json:"checkout" form:"checkout"`
}
