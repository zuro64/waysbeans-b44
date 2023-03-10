package routes

import (
	"nis-waybeans/handlers"
	"nis-waybeans/pkg/middleware"
	"nis-waybeans/pkg/mysql"
	"nis-waybeans/repositories"

	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Group) {
	productRepository := repositories.RepositoryProduct(mysql.DB)
	cartRepository := repositories.RepositoryCart(mysql.DB)
	h := handlers.HandlerProduct(
		productRepository,
		cartRepository,
	)

	e.GET("/products", h.FindProduct)
	e.GET("/product/:id", h.GetProduct)
	e.GET("/products/search", h.SearchProduct)
	e.GET("/products/top", h.FindTopProducts)
	e.POST("/product", middleware.Auth(middleware.UploadFile(h.CreateProduct, "image")))
	e.PATCH("/product/:id", middleware.Auth(middleware.UploadFile(h.UpdateProduct, "image")))
	e.DELETE("/product/:id", middleware.Auth(h.DeleteProduct))

}
