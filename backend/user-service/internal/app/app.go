package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunServe() {
	e := echo.New()
	e.Use(middleware.CORS())

	customValidator :=
}