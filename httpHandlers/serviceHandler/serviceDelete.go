package serviceHandler

import (
	"errors"
	"net/http"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type ServiceDeleteRequest struct {
	ID int64 `json:"id"`
}

func ServiceDelete(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ServiceDeleteRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ServiceDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result := db.Delete(&dataLayer.Service{}, request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with delete service by id on ServiceDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}
