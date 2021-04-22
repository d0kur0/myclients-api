package userHandler

import (
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserSignUpRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

type UserSignUpResponse struct {
	IsError bool           `json:"isError"`
	Errors  []string       `json:"errors"`
	User    dataLayer.User `json:"user"`
	Token   string         `json:"token"`
}

func UserSignUp(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(UserSignUpRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on UserSignUp")
		return c.JSON(http.StatusInternalServerError, "")
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on UserSignUp")
		return c.JSON(http.StatusInternalServerError, "")
	}

	if validateErrs != nil {
		return c.JSON(http.StatusOK, UserSignUpResponse{
			IsError: true,
			Errors:  validateErrs,
		})
	}

	var foundEmails int64 = 0
	err = db.Model(&dataLayer.User{}).Select("email").Where("email = ?", request.Email).Count(&foundEmails).Error
	if err != nil {
		c.Logger().Error("Error with found email on UserSignUp")
		return c.JSON(http.StatusInternalServerError, "")
	}

	if foundEmails > 0 {
		return c.JSON(http.StatusOK, UserSignUpResponse{
			IsError: true,
			Errors:  []string{"This email address is already in use"},
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error("Error with hashing password on UserSignUp")
		return c.JSON(http.StatusInternalServerError, "")
	}

	user := dataLayer.User{
		Name:           request.Name,
		Email:          request.Email,
		Password:       string(hashedPassword),
		IsEmailConfirm: false,
	}

	if err := db.Create(&user).Error; err != nil {
		c.Logger().Error("Error save user in db on UserSignUp")
		return c.JSON(http.StatusInternalServerError, "")
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

	return c.JSON(http.StatusOK, UserSignUpResponse{
		IsError: false,
		User:    user,
		Token:   newAuthToken.Token,
	})
}
