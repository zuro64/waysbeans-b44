package routes

import (
	"nis-waybeans/handlers"
	"nis-waybeans/pkg/middleware"
	"nis-waybeans/pkg/mysql"
	"nis-waybeans/repositories"

	"github.com/labstack/echo/v4"
)

func CartRoutes(e *echo.Group) {
	transactionRepository := repositories.RepositoryTransaction(mysql.DB)
	productRepository := repositories.RepositoryProduct(mysql.DB)
	cartRepository := repositories.RepositoryCart(mysql.DB)

	h := handlers.HandlerCart(
		transactionRepository,
		productRepository,
		cartRepository,
	)

	e.GET("/carts", h.FindCarts)
	e.GET("/carts/not-checkout", middleware.Auth(h.FindUnCheckedOutCarts))
	e.GET("/cart/:id", h.GetCart)
	e.POST("/cart", middleware.Auth(h.CreateCart))
	e.PATCH("/cart/:id", middleware.Auth(h.UpdateCart))
	e.DELETE("/cart/:id", h.DeleteCart)

}
