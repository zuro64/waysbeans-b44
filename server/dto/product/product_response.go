package productdto

import "time"

type ProductResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" form:"name" validate:"required"`
	Price       int       `json:"price" form:"price" validate:"required"`
	Description string    `json:"description" form:"description" validate:"required"`
	Stock       int       `json:"stock" form:"stock" validate:"required"`
	Image       string    `json:"image" form:"image" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
