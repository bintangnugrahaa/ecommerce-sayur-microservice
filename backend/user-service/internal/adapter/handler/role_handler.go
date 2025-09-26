package handler

import (
	"user-service/config"
	"user-service/internal/core/service"

	"github.com/labstack/echo/v4"
)

type RoleHandlerInterface interface {
	GetAll(c echo.Context) error
	GetByID(c echo.Context) error
	Create(c echo.Context) error
	Delete(c echo.Context) error
	Update(c echo.Context) error
}

type roleHandler struct {
	roleService service.RoleServiceInterface
}

func NewRoleHandler(e *echo.Echo, roleService service.RoleServiceInterface, cfg *config.Config, jwtService service.JwtServiceInterface) RoleHandlerInterface {
	roleHandlers := &roleHandler{}
}
