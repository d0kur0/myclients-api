package httpHandlers

import (
	"github.com/d0kur0/myclients-api/httpHandlers/clientHandler"
	"github.com/d0kur0/myclients-api/httpHandlers/recordHandler"
	"github.com/d0kur0/myclients-api/httpHandlers/userHandler"
	"github.com/labstack/echo"
)

type Route struct {
	Handler    func(c echo.Context) (err error)
	IsNeedAuth bool
	Path       string
}

var routes = []Route{
	// User actions
	{
		Path:       "/user/signIn",
		Handler:    userHandler.UserSignIn,
		IsNeedAuth: false,
	},
	{
		Path:       "/user/signUp",
		Handler:    userHandler.UserSignUp,
		IsNeedAuth: false,
	},
	{
		Path:       "/user/resetPassword",
		Handler:    userHandler.UserResetPassword,
		IsNeedAuth: false,
	},
	// Client actions
	{
		Path:       "/client/getAll",
		Handler:    clientHandler.ClientGetAll,
		IsNeedAuth: true,
	},
	{
		Path:       "/client/create",
		Handler:    clientHandler.ClientCreate,
		IsNeedAuth: true,
	},
	{
		Path:       "/client/update",
		Handler:    clientHandler.ClientUpdate,
		IsNeedAuth: true,
	},
	{
		Path:       "/client/delete",
		Handler:    clientHandler.ClientDelete,
		IsNeedAuth: true,
	},
	// Service actions
	{
		Path:       "/service/getAll",
		Handler:    ServiceGetAll,
		IsNeedAuth: true,
	},
	{
		Path:       "/service/create",
		Handler:    ServiceCreate,
		IsNeedAuth: true,
	},
	{
		Path:       "/service/update",
		Handler:    ServiceUpdate,
		IsNeedAuth: true,
	},
	{
		Path:       "/service/delete",
		Handler:    ServiceDelete,
		IsNeedAuth: true,
	},
	// Record actions
	{
		Path:       "/record/getAll",
		Handler:    recordHandler.RecordGetAll,
		IsNeedAuth: true,
	},
	{
		Path:       "/record/getByDate",
		Handler:    recordHandler.RecordGetByDate,
		IsNeedAuth: true,
	},
	{
		Path:       "/record/create",
		Handler:    recordHandler.RecordCreate,
		IsNeedAuth: true,
	},
	{
		Path:       "/record/update",
		Handler:    recordHandler.RecordUpdate,
		IsNeedAuth: true,
	},
	{
		Path:       "/record/delete",
		Handler:    recordHandler.RecordDelete,
		IsNeedAuth: true,
	},
}
