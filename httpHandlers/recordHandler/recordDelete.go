package recordHandler

import (
	"errors"
	"net/http"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type RecordDeleteRequest struct {
	ID uint64 `json:"id"`
}

func RecordDelete(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(RecordDeleteRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordDelete", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	result := db.Delete(&dataLayer.Record{}, request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with delete client by id on RecordDelete", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, "")
}
