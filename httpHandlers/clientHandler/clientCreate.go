package clientHandler

import (
	"net/http"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
)

type ClientCreateRequest struct {
	FirstName   string `json:"firstName" validate:"required,max=32"`
	MiddleName  string `json:"middleName" validate:"max=32"`
	Description string `json:"description" validate:"max=256"`
}

type ClientCreateResponse struct {
	IsError bool             `json:"isError"`
	Errors  []string         `json:"errors"`
	Client  dataLayer.Client `json:"client"`
}

func ClientCreate(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(ClientCreateRequest)
	requestUser, err := helpers.GetUserByRequest(c)

	if err != nil {
		c.Logger().Error("Request user not found", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on ClientCreate", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on ClientCreate", err)
		return c.JSON(http.StatusOK, ClientCreateResponse{
			IsError: true,
			Errors:  []string{"internal server error"},
		})
	}

	if validateErrs != nil {
		return c.JSON(http.StatusOK, ClientCreateResponse{
			IsError: true,
			Errors:  validateErrs,
		})
	}

	client := dataLayer.Client{
		FirstName:   request.FirstName,
		MiddleName:  request.MiddleName,
		Description: request.Description,
		UserID:      requestUser.ID,
	}

	if err := db.Create(&client).Error; err != nil {
		c.Logger().Error("Error save user in db on ClientCreate")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, ClientCreateResponse{
		IsError: false,
		Errors:  nil,
		Client:  client,
	})
}
