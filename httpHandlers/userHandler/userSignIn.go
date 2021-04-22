package userHandler

import (
	"errors"
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/labstack/echo"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignInResponse struct {
	IsError bool           `json:"isError"`
	Errors  []string       `json:"errors"`
	Token   string         `json:"token"`
	User    dataLayer.User `json:"user"`
}

func UserSignIn(c echo.Context) (err error) {
	request := new(UserSignInRequest)
	db := dataLayer.GetDB()

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on UserSignIn")
		return c.JSON(http.StatusInternalServerError, "")
	}

	var user dataLayer.User
	result := db.Where(&dataLayer.User{Email: request.Email}).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusOK, UserSignInResponse{
			IsError: true,
			Errors:  []string{"This email not found"},
		})
	}

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with found user by email on UserSignIn")
		return c.JSON(http.StatusInternalServerError, "")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusOK, UserSignInResponse{
			IsError: true,
			Errors:  []string{"Password incorrect"},
		})
	}

	newAuthToken := dataLayer.AuthToken{
		UserID:   user.ID,
		Token:    randstr.Hex(32),
		DeadTime: time.Now().AddDate(0, 6, 0),
	}

	createResult := db.Create(&newAuthToken)
	if createResult.Error != nil {
		c.Logger().Error("Error with create user AuthToken on UserSignIn")
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, UserSignInResponse{
		IsError: false,
		User:    user,
		Token:   newAuthToken.Token,
	})
}
