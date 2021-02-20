package httpHandlers

import (
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
)

type UserSignUpRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

func UserSignUp(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(UserSignUpRequest)

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on UserSignUp")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	validateErrs, err := helpers.Validate(request)
	if err != nil {
		c.Logger().Error("Error with validate request struct on UserSignUp")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if validateErrs != nil {
		return c.JSON(http.StatusBadRequest, validateErrs)
	}

	var foundEmails int64 = 0
	err = db.Model(&dataLayer.User{}).Select("email").Where("email = ?", request.Email).Count(&foundEmails).Error
	if err != nil {
		c.Logger().Error("Error with found email on UserSignUp")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if foundEmails > 0 {
		return c.JSON(http.StatusBadRequest, []string{"This email address is already in use"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Error("Error with hashing password on UserSignUp")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	user := dataLayer.User{
		Name:           request.Name,
		Email:          request.Email,
		Password:       string(hashedPassword),
		IsEmailConfirm: false,
	}

	if err := db.Create(&user).Error; err != nil {
		c.Logger().Error("Error save user in db on UserSignUp")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, "")
}

type UserSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignInResponse struct {
	Token string `json:"token"`
}

func UserSignIn(c echo.Context) (err error) {
	request := new(UserSignInRequest)
	db := dataLayer.GetDB()

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on UserSignIn")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	var user dataLayer.User
	result := db.Where(&dataLayer.User{Email: request.Email}).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusBadRequest, []string{"This email not found"})
	}

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Logger().Error("Error with found user by email on UserSignIn")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, []string{"Password incorrect"})
	}

	newAuthToken := dataLayer.AuthToken{
		UserID:   user.ID,
		Token:    randstr.Hex(32),
		LifeTime: time.Now().AddDate(0, 6, 0),
	}

	createResult := db.Create(&newAuthToken)
	if createResult.Error != nil {
		c.Logger().Error("Error with create user AuthToken on UserSignIn")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, UserSignInResponse{newAuthToken.Token})
}

func UserResetPassword(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}
