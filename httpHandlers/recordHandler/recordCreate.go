package recordHandler

import (
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
)

type RecordCreateRequest struct {
	ServiceIDs []uint64  `json:"serviceIds" validate:"required"`
	ClientID   uint64    `json:"clientId" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
}

type RecordCreateResponse struct {
	IsError bool             `json:"isError"`
	Errors  []string         `json:"errors"`
	Record  dataLayer.Record `json:"record"`
}

func RecordCreate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(RecordCreateRequest)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordCreate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on RecordCreate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	if validateErrs != nil {
		return c.JSON(http.StatusOK, RecordCreateResponse{
			IsError: true,
			Errors:  validateErrs,
		})
	}

	var services []dataLayer.Service
	result := db.Find(&services, request.ServiceIDs)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find all services by ids on RecordCreate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	record := dataLayer.Record{
		UserID:   requestUser.ID,
		ClientID: request.ClientID,
		Date:     request.Date,
	}

	if err := db.Create(&record).Error; err != nil {
		c.Logger().Error("Error save user in db on RecordCreate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	err = db.Model(&record).Association("Services").Append(&services)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, RecordCreateResponse{
		IsError: false,
		Record:  record,
	})
}
