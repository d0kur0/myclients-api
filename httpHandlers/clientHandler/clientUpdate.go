package clientHandler

import (
	"errors"
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

type ClientUpdateRequest struct {
	ID          int64  `json:"id" validate:"required"`
	FirstName   string `json:"firstName" validate:"required,max=32"`
	MiddleName  string `json:"middleName" validate:"required,max=32"`
	Description string `json:"description" validate:"required,max=256"`
}

type ClientUpdateResponse struct {
	IsError bool             `json:"isError"`
	Errors  []string         `json:"errors"`
	Client  dataLayer.Client `json:"client"`
}

func ClientUpdate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ClientUpdateRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ClientUpdate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on ClientUpdate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	if validateErrs != nil {
		return c.JSON(http.StatusOK, ClientUpdateResponse{
			IsError: true,
			Errors:  validateErrs,
		})
	}

	var client dataLayer.Client

	result := db.Model(&client).Where("id = ?", request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find client by id on ClientUpdate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	result = db.Model(&client).Where("id = ?", request.ID).Updates(&dataLayer.Client{
		FirstName:   request.FirstName,
		MiddleName:  request.MiddleName,
		Description: request.Description,
	})

	if result.Error != nil {
		c.Logger().Error("Error with update client by id on ClientUpdate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, ClientUpdateResponse{
		IsError: false,
		Errors:  nil,
		Client:  client,
	})
}
