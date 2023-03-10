package transactiondto

type CreateTransactionRequest struct {
	UserID     int    `json:"user_id" form:"user_id"`
	Name       string `json:"name" form:"name"`
	Email      string `json:"email"  form:"email"`
	Phone      string `json:"phone"  form:"phone"`
	PostCode   string `json:"post_code"  form:"post_code"`
	Address    string `json:"address"  form:"address"`
	Attachment string `json:"attachment"  form:"attachment"`
	Status     string `json:"status"  form:"status"`
}

type UpdateTransactionRequest struct {
	Name       string `json:"name" form:"name" validate:"required"`
	Email      string `json:"email"  form:"email" validate:"required"`
	Phone      string `json:"phone"  form:"phone" validate:"required"`
	PostCode   string `json:"post_code"  form:"post_code" validate:"required"`
	Address    string `json:"address"  form:"address" validate:"required"`
	Attachment string `json:"attachment"  form:"attachment" validate:"required"`
	Status     string `json:"status"  form:"status" validate:"required"`
	TotalQty   int    `json:"total_qty" gorm:"type: int" form:"total_qty"`
	SubTotal   int    `json:"sub_total" gorm:"type: int" form:"sub_total"`
}
