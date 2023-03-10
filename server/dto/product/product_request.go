package productdto

type CreateProductRequest struct {
	Name        string `json:"name" form:"name"  validate:"required"`
	Price       int    `json:"price" form:"price" validate:"required"`
	Description string `json:"description" form:"description"  validate:"required"`
	Stock       int    `json:"stock" form:"stock" validate:"required"`
	Image       string `json:"image" form:"image" validate:"required"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" form:"name"`
	Price       int    `json:"price" form:"price"`
	Description string `json:"description" form:"description"`
	Stock       int    `json:"stock" form:"stock"`
	Image       string `json:"image" form:"image"`
}
