package serviceHandler

import (
	"net/http"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
)

type ServiceCreateRequest struct {
	Name  string `json:"name" validate:"required,max=64"`
	Price int64  `json:"price" validate:"required"`
}

type ServiceCreateResponse struct {
	IsError bool              `json:"isError"`
	Errors  []string          `json:"errors"`
	Service dataLayer.Service `json:"service"`
}

func ServiceCreate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ServiceCreateRequest)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ServiceCreate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on ServiceCreate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	if validateErrs != nil {
		return c.JSON(http.StatusOK, ServiceCreateResponse{
			IsError: true,
			Errors:  validateErrs,
		})
	}

	service := dataLayer.Service{
		Name:   request.Name,
		Price:  request.Price,
		UserID: requestUser.ID,
	}

	if err := db.Create(&service).Error; err != nil {
		c.Logger().Error("Error save user in db on ServiceCreate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, ServiceCreateResponse{
		IsError: false,
		Service: service,
	})
}
