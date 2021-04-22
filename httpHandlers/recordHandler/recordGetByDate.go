package recordHandler

import (
	"errors"
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type RecordGetByDateRequest struct {
	Day   int
	Month time.Month
	Year  int
}

func RecordGetByDate(c echo.Context) (err error) {
	db := dataLayer.GetDB()

	request := new(RecordGetByDateRequest)
	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordGetByDate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	requestUser, err := helpers.GetUserByRequest(c)
	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	requestDate := time.Date(request.Year, request.Month, request.Day, 0, 0, 0, 0, time.Local)
	records := new([]dataLayer.Record)

	result := db.Preload("Service").Preload("Client").Find(
		&records,
		"user_id = ? AND DATE(date) = ?",
		requestUser.ID, requestDate.Format("2006-01-02"),
	)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find records by date in db", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, records)
}
