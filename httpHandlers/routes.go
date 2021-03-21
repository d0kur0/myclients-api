package httpHandlers

import "github.com/labstack/echo"

type Route struct {
	Handler    func(c echo.Context) (err error)
	IsNeedAuth bool
	Path       string
}

var routes = []Route{
	// User actions
	{
		Path:       "/user/signIn",
		Handler:    UserSignIn,
		IsNeedAuth: false,
	},
	{
		Path:       "/user/signUp",
		Handler:    UserSignUp,
		IsNeedAuth: false,
	},
	{
		Path:       "/user/resetPassword",
		Handler:    UserResetPassword,
		IsNeedAuth: false,
	},
	// Client actions
	{
		Path:       "/client/getAll",
		Handler:    ClientGetAll,
		IsNeedAuth: true,
	},
	{
		Path:       "/client/create",
		Handler:    ClientCreate,
		IsNeedAuth: true,
	},
	{
		Path:       "/client/update",
		Handler:    ClientUpdate,
		IsNeedAuth: true,
	},
	{
		Path:       "/client/delete",
		Handler:    ClientDelete,
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
		Handler:    RecordGetAll,
		IsNeedAuth: true,
	},
	{
		Path:       "/record/getByDate",
		Handler:    RecordGetByDate,
		IsNeedAuth: true,
	},
	{
		Path:       "/record/create",
		Handler:    RecordCreate,
		IsNeedAuth: true,
	},
	{
		Path:       "/record/update",
		Handler:    RecordUpdate,
		IsNeedAuth: true,
	},
	{
		Path:       "/record/delete",
		Handler:    RecordDelete,
		IsNeedAuth: true,
	},
}
