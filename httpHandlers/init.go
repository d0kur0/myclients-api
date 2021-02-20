package httpHandlers

import (
	"errors"
	"net/http"
	"reflect"
	"time"

	"gorm.io/gorm"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
)

func entryPoint(c echo.Context) (err error) {
	var foundedRoute Route
	for _, route := range routes {
		if route.Path == c.Request().RequestURI {
			foundedRoute = route
			break
		}
	}

	if reflect.ValueOf(foundedRoute).IsZero() {
		return c.JSON(http.StatusNotFound, "")
	}

	if foundedRoute.IsNeedAuth {
		var authToken dataLayer.AuthToken

		db := dataLayer.GetDB()
		result := db.Where("token = ?", c.Request().Header.Get("AuthToken")).First(&authToken)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			c.Logger().Error("Error with found auth token in database on EntryPoint")
			return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		}

		if time.Now().After(authToken.LifeTime) {
			deleteResult := db.Delete(&authToken)

			if deleteResult.Error != nil && errors.Is(deleteResult.Error, gorm.ErrRecordNotFound) {
				c.Logger().Error("Error with delete authToken on EntryPoint")
				return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
			}

			return c.JSON(http.StatusUnauthorized, []string{"auth token is dead"})
		}

		if reflect.ValueOf(authToken).IsZero() {
			return c.JSON(http.StatusUnauthorized, []string{"Method required authorization"})
		}
	}

	return foundedRoute.Handler(c)
}

func Init(e *echo.Echo) (err error) {
	e.Any("/*", entryPoint)

	return err
}
