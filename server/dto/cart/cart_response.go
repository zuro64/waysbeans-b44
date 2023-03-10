package cartdto

import "time"

type CartResponse struct {
	ID            int       `json:"id"`
	ProductID     int       `json:"product_id" form:"product_id" validate:"required"`
	OrderQuantity int       `json:"order_quantity" form:"order_quantity" validate:"required"`
	Checkout      bool      `json:"checkout" form:"checkout" validate:"required"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
