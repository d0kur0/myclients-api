package httpHandlers

import (
	"github.com/labstack/echo"
)

type route struct {
	Handler    func(c echo.Context) (err error)
	IsNeedAuth bool
	Path       string
}

var routes = []route{
	{
		Path:       "/user/signIn",
		Handler:    UserSignIn,
		IsNeedAuth: false,
	},
	{
		Path: "/user/signUp",
	},
}

func entryPoint(c echo.Context) (err error) {
	return nil
}

func Init(e *echo.Echo) (err error) {
	e.Any("*", entryPoint)

	return err
}
