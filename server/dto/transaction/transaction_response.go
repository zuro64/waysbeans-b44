package transactiondto

import (
	"nis-waybeans/models"
	"time"
)

type TransactionResponse struct {
	ID         int                 `json:"id"`
	User       models.UserResponse `json:"user"`
	Name       string              `json:"name"  form:"name" validate:"required"`
	Email      string              `json:"email"  form:"email" validate:"required"`
	Phone      string              `json:"phone"  form:"phone" validate:"required"`
	Address    string              `json:"address"  form:"address" validate:"required"`
	Attachment string              `json:"attachment"  form:"attachment" validate:"required"`
	Status     string              `json:"status"  form:"attachment" validate:"required"`
	SubTotal   int                 `json:"sub_total"`
	TotalQty   int                 `json:"total_qty"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`

	Cart []models.CartResponse `json:"products"`
}
