package recordHandler

import (
	"errors"
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

func RecordGetAll(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	records := new([]dataLayer.Record)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result := db.Preload("Service").Preload("Client").Find(&records, "user_id = ?", requestUser.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find all records in db", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, records)
}
