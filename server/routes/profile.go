package routes

import (
	"nis-waybeans/handlers"
	"nis-waybeans/pkg/middleware"
	"nis-waybeans/pkg/mysql"
	"nis-waybeans/repositories"

	"github.com/labstack/echo/v4"
)

func ProfileRoutes(e *echo.Group) {
	ProfileRepository := repositories.RepositoryProfile(mysql.DB)
	h := handlers.HandlerProfile(ProfileRepository)

	e.GET("/profiles", h.FindProfile)
	e.GET("/profile/:id", h.GetProfile)
	e.POST("/profile", middleware.Auth(middleware.UploadFile(h.CreateProfile, "profile_picture")))
	e.PATCH("/profile", middleware.Auth(middleware.UploadFile(h.UpdateProfile, "profile_picture")))
	e.DELETE("/profile/:id", h.DeleteProfile)
}
