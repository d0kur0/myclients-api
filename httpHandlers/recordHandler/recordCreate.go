package recordHandler

import (
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

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
