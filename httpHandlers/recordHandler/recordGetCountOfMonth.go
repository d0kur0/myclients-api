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

type RecordGetCountOfMonthRequest struct {
	Date string `json:"date"`
}

type CountOfDay struct {
	Day   int `json:"day"`
	Count int `json:"count"`
}

type RecordGetCountOfMonthResponse = []CountOfDay

func RecordGetCountOfMonth(c echo.Context) (err error) {
	db := dataLayer.GetDB()

	request := new(RecordGetCountOfMonthRequest)
	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordGetCountOfMonth", err)
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

	var records []dataLayer.Record
	result := db.Find(
		&records,
		"user_id = ? AND DATE(date, 'localtime') LIKE ?",
		requestUser.ID, fmt.Sprintf("%d-%s-%%", requestDate.Year(), fmt.Sprintf("%02d", requestDate.Month())),
	)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find records by date in db", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	countOfDays := make(map[int]int)

	for _, record := range records {
		if _, exists := countOfDays[record.Date.Day()]; exists {
			countOfDays[record.Date.Day()] = countOfDays[record.Date.Day()] + 1
		} else {
			countOfDays[record.Date.Day()] = 1
		}
	}

	var response RecordGetCountOfMonthResponse
	for day, count := range countOfDays {
		response = append(response, CountOfDay{
			Day:   day,
			Count: count,
		})
	}

	return c.JSON(http.StatusOK, response)
}
