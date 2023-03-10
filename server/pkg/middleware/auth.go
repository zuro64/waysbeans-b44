package middleware

import (
	"net/http"
	dto "nis-waybeans/dto/result"
	jwtToken "nis-waybeans/pkg/jwt"
	"strings"

	"github.com/labstack/echo/v4"
)

// Result struct
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Auth function
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResult{
				Code:    http.StatusBadRequest,
				Message: "unauthrized",
			})
		}

		token = strings.Split(token, " ")[1]
		claims, err := jwtToken.DecodeToken(token)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, Result{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
		}
		c.Set("userLogin", claims)
		return next(c)
	}
}
