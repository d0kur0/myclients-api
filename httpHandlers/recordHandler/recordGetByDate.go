package recordHandler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type RecordGetByDateRequest struct {
	Date string `json:"date"`
}

func RecordGetByDate(c echo.Context) (err error) {
	db := dataLayer.GetDB()

	request := new(RecordGetByDateRequest)
	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordGetByDate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	requestUser, err := helpers.GetUserByRequest(c)
	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	requestDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error on parse requestDate (by format yyyy-MM-dd): %s", err))
	}

	records := new([]dataLayer.Record)

	result := db.Preload("Services").Preload("Client").Find(
		&records,
		"user_id = ? AND DATE(date, 'localtime') = ?",
		requestUser.ID, requestDate.Format("2006-01-02"),
	)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find records by date in db", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, records)
}
