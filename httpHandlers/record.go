package httpHandlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/d0kur0/myclients-api/helpers"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
	"gorm.io/gorm"
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

type RecordCreateRequest struct {
	ServiceID uint64    `json:"serviceId" validate:"required"`
	ClientID  uint64    `json:"clientId" validate:"required"`
	Date      time.Time `json:"date" validate:"required"`
}

func RecordCreate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(RecordCreateRequest)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordCreate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on RecordCreate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if validateErrs != nil {
		return c.JSON(http.StatusBadRequest, validateErrs)
	}

	record := dataLayer.Record{
		UserID:    requestUser.ID,
		ServiceID: request.ServiceID,
		ClientID:  request.ClientID,
		Date:      request.Date,
	}

	if err := db.Create(&record).Error; err != nil {
		c.Logger().Error("Error save user in db on RecordCreate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}

type RecordUpdateRequest struct {
	ID     int64               `json:"id" validate:"required"`
	Record RecordCreateRequest `json:"client" validate:"required,dive"`
}

func RecordUpdate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(RecordUpdateRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordUpdate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on RecordUpdate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if validateErrs != nil {
		return c.JSON(http.StatusBadRequest, validateErrs)
	}

	var record dataLayer.Record

	result := db.Model(&record).Where("id = ?", request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find client by id on RecordUpdate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result = db.Model(&record).Where("id = ?", request.ID).Updates(&dataLayer.Record{
		ServiceID: request.Record.ServiceID,
		ClientID:  request.Record.ClientID,
		Date:      request.Record.Date,
	})

	if result.Error != nil {
		c.Logger().Error("Error with update client by id on RecordUpdate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}

type RecordDeleteRequest struct {
	ID uint64 `json:"id"`
}

func RecordDelete(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(RecordDeleteRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result := db.Delete(&dataLayer.Record{}, request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with delete client by id on RecordDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}
