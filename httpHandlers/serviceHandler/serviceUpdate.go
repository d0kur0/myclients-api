package serviceHandler

import (
	"errors"
	"net/http"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type ServiceUpdateRequest struct {
	ID    int64  `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,max=64"`
	Price int64  `json:"price" validate:"required"`
}

type ServiceUpdateResponse struct {
	IsError bool              `json:"isError"`
	Errors  []string          `json:"errors"`
	Service dataLayer.Service `json:"service"`
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
		return c.JSON(http.StatusOK, ServiceUpdateResponse{
			IsError: true,
			Errors:  validateErrs,
		})
	}

	var service dataLayer.Service

	result := db.Model(&service).Where("id = ?", request.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with find service by id on ServiceUpdate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	result = db.Model(&service).Where("id = ?", request.ID).Updates(&dataLayer.Service{
		Name:  request.Name,
		Price: request.Price,
	})

	if result.Error != nil {
		c.Logger().Error("Error with update service by id on ServiceUpdate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, ServiceUpdateResponse{
		IsError: false,
		Service: service,
	})
}
