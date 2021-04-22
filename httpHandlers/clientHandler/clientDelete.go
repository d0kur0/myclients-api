package clientHandler

import (
	"errors"
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

type ClientDeleteRequest struct {
	ID uint64 `json:"id"`
}

func ClientDelete(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ClientDeleteRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ClientDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result := db.Delete(&dataLayer.Client{}, request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with delete client by id on ClientDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}
