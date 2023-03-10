package middleware

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc, imageKey string) echo.HandlerFunc {
	return func(c echo.Context) error {

		file, err := c.FormFile(imageKey)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		defer src.Close()

		tempFile, err := ioutil.TempFile("uploads", "image-*.png")
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		defer tempFile.Close()

		if _, err = io.Copy(tempFile, src); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		data := tempFile.Name()
		filename := data[8:] // split uploads/

		c.Set("dataFile", filename)
		return next(c)
	}
}
