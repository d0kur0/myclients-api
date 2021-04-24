package recordHandler

import (
	"errors"
	"net/http"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func RecordGetAll(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	records := new([]dataLayer.Record)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	result := db.Preload("Services").Preload("Client").Find(&records, "user_id = ?", requestUser.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find all records in db", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, records)
}
