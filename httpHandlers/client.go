package httpHandlers

import (
	"errors"
	"net/http"

	"github.com/d0kur0/myclients-api/helpers"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func ClientGetAll(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	clients := new([]dataLayer.Client)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result := db.Find(&clients, "user_id = ?", requestUser.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find all clients in db")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, clients)
}

type ClientCreateRequest struct {
	FirstName   string `json:"firstName" validate:"required,max=32"`
	MiddleName  string `json:"middleName" validate:"required,max=32"`
	Description string `json:"description" validate:"required,max=256"`
}

func ClientCreate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ClientCreateRequest)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ClientCreate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on ClientCreate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if validateErrs != nil {
		return c.JSON(http.StatusBadRequest, validateErrs)
	}

	client := dataLayer.Client{
		FirstName:   request.FirstName,
		MiddleName:  request.MiddleName,
		Description: request.Description,
		UserID:      requestUser.ID,
	}

	if err := db.Create(&client).Error; err != nil {
		c.Logger().Error("Error save user in db on ClientCreate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}

type ClientUpdateRequest struct {
	ID     int64               `json:"id" validate:"required"`
	Client ClientCreateRequest `json:"client" validate:"required,dive"`
}

func ClientUpdate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ClientUpdateRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ClientUpdate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on ClientUpdate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if validateErrs != nil {
		return c.JSON(http.StatusBadRequest, validateErrs)
	}

	var client dataLayer.Client

	result := db.Model(&client).Where("id = ?", request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find client by id on ClientUpdate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result = db.Model(&client).Where("id = ?", request.ID).Updates(&dataLayer.Client{
		FirstName:   request.Client.FirstName,
		MiddleName:  request.Client.MiddleName,
		Description: request.Client.Description,
	})

	if result.Error != nil {
		c.Logger().Error("Error with update client by id on ClientUpdate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}

type ClientDeleteRequest struct {
	ID uint64 `json:"id"`
}

func ClientDelete(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ClientDeleteRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ClientDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result := db.Delete(&dataLayer.Client{}, request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with delete client by id on ClientDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}
