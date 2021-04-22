package userHandler

import (
	"github.com/labstack/echo"
	"net/http"
)

func UserResetPassword(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}
