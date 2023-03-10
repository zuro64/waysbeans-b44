package models

import "time"

type Product struct {
	ID          int                 `json:"id" gorm:"primary_key:auto_increment"`
	Name        string              `json:"name" form:"name" gorm:"type: varchar(255)"`
	Price       int                 `json:"price" form:"price" gorm:"type: int"`
	Description string              `json:"description" form:"description" gorm:"type: varchar(255)"`
	Stock       int                 `json:"stock" form:"stock" gorm:"type: int"`
	Image       string              `json:"image" form:"image" gorm:"type: varchar(255)"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Cart        CartProductResponse `json:"-"`
}

type ProductResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func (ProductResponse) TableName() string {
	return "products"
}
