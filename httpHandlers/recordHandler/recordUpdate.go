package recordHandler

import (
	"errors"
	"net/http"
	"time"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type RecordUpdateRequest struct {
	ID         int64     `json:"id" validate:"required"`
	ServiceIDs []uint64  `json:"serviceIds" validate:"required"`
	ClientID   uint64    `json:"clientId" validate:"required"`
	Date       time.Time `json:"date" validate:"required"`
}

type RecordUpdateResponse struct {
	IsError bool             `json:"isError"`
	Errors  []string         `json:"errors"`
	Record  dataLayer.Record `json:"record"`
}

func RecordUpdate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(RecordUpdateRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on RecordUpdate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on RecordUpdate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	if validateErrs != nil {
		return c.JSON(http.StatusOK, RecordUpdateResponse{
			IsError: true,
			Errors:  validateErrs,
		})
	}

	var record dataLayer.Record

	result := db.Model(&record).Where("id = ?", request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find client by id on RecordUpdate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	var services []dataLayer.Service
	for _, serviceId := range request.ServiceIDs {
		services = append(services, dataLayer.Service{
			UserID: serviceId,
		})
	}

	result = db.Model(&record).Where("id = ?", request.ID).Updates(&dataLayer.Record{
		Services: services,
		ClientID: request.ClientID,
		Date:     request.Date,
	})

	if result.Error != nil {
		c.Logger().Error("Error with update client by id on RecordUpdate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, RecordUpdateResponse{
		IsError: false,
		Record:  record,
	})
}
