package httpHandlers

import (
	"errors"
	"net/http"

	"github.com/d0kur0/myclients-api/helpers"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func ServiceGetAll(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	services := new([]dataLayer.Service)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result := db.Find(&services, "user_id = ?", requestUser.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find all services in db", result.Error)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, services)
}

type ServiceCreateRequest struct {
	Name  string `json:"name" validate:"required,max=32"`
	Price int64  `json:"price" validate:"required"`
}

func ServiceCreate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ServiceCreateRequest)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ServiceCreate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on ServiceCreate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if validateErrs != nil {
		return c.JSON(http.StatusBadRequest, validateErrs)
	}

	service := dataLayer.Service{
		Name:   request.Name,
		Price:  request.Price,
		UserID: requestUser.ID,
	}

	if err := db.Create(&service).Error; err != nil {
		c.Logger().Error("Error save user in db on ServiceCreate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}

type ServiceUpdateRequest struct {
	ID      int64                `json:"id" validate:"required"`
	Service ServiceCreateRequest `json:"service" validate:"required,dive"`
}

func ServiceUpdate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ServiceUpdateRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ServiceUpdate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on ServiceUpdate", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if validateErrs != nil {
		return c.JSON(http.StatusBadRequest, validateErrs)
	}

	var service dataLayer.Service

	result := db.Model(&service).Where("id = ?", request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find service by id on ServiceUpdate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result = db.Model(&service).Where("id = ?", request.ID).Updates(&dataLayer.Service{
		Name:  request.Service.Name,
		Price: request.Service.Price,
	})

	if result.Error != nil {
		c.Logger().Error("Error with update service by id on ServiceUpdate")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}

type ServiceDeleteRequest struct {
	ID int64 `json:"id"`
}

func ServiceDelete(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ServiceDeleteRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ServiceDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	result := db.Delete(&dataLayer.Service{}, request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with delete service by id on ServiceDelete", err)
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}
