package serviceHandler

import (
	"errors"
	"net/http"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
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
