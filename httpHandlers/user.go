package httpHandlers

import (
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/helpers"
	"github.com/labstack/echo"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

func UserSignUp(c echo.Context) (err error) {
	db := dataLayer.GetDB()
	request := new(SignUpRequest)

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

	// You are gay?
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

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserSignIn(c echo.Context) (err error) {
	request := new(SignInRequest)
	db := dataLayer.GetDB()

	if err = c.Bind(request); err != nil {
		c.Logger().Error("Bind struct error on UserSignIn")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	var user dataLayer.User
	result := db.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		c.Logger().Error("Error with found user by email on UserSignIn")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, []string{"This email not found"})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, []string{"Password incorrect"})
	}

	newAuthToken := dataLayer.AuthToken{UserID: user.ID, Token: randstr.Hex(32)}
	createResult := db.Create(&newAuthToken)
	if createResult.Error != nil {
		c.Logger().Error("Error with create user AuthToken on UserSignIn")
		return c.JSON(http.StatusInternalServerError, []string{"internal server error"})
	}

	return c.JSON(http.StatusOK, newAuthToken.Token)
}

func UserResetPassword(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}
