package helpers

import (
	"errors"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetUserByRequest(c echo.Context) (user dataLayer.User, err error) {
	db := dataLayer.GetDB()

	var authToken dataLayer.AuthToken
	tokenResult := db.Where("token = ?", c.Request().Header.Get("AuthToken")).First(&authToken)
	if tokenResult.Error != nil && !errors.Is(tokenResult.Error, gorm.ErrRecordNotFound) {
		err = errors.New("error with found auth token in database on GetUserByRequest")
		return
	}

	userResult := db.Where("id = ?", authToken.UserID).First(&user)
	if userResult.Error != nil && !errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
		err = errors.New("error with found user in database on GetUserByRequest")
		return
	}

	return
}
