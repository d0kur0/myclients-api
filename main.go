package main

import (
	"fmt"
	"github.com/d0kur0/myclients-api/dataLayer"
	"github.com/d0kur0/myclients-api/httpHandlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"os"
)

func main() {
	e := echo.New()

	if err := godotenv.Load(".env"); err != nil {
		e.Logger.Fatal("Error loading .env file")
		return
	}

	if err := dataLayer.Init(); err != nil {
		e.Logger.Fatal("Error init dataLayer")
		return
	}

	if err := httpHandlers.Init(e); err != nil {
		e.Logger.Fatal("Error init httpHandlers")
		return
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))))
}
