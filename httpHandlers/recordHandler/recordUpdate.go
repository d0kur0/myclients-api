package recordHandler

import (
	"errors"
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

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
