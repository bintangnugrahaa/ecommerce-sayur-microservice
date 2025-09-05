package handler

import (
	"net/http"
	"user-service/internal/adapter/handler/request"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/service"

	"github.com/google/martian/v3/log"
	"github.com/labstack/echo/v4"
)

type UserHandlerInterface interface {
	SignIn(c echo.Context) error
}

type userHandler struct {
	userService service.UserServiceInterface
}

var err error

func (u *userHandler) SignIn(c echo.Context) error {
	var (
		req        = request.SignInRequest{}
		resp       = response.DefaultResponse{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	if err = c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-1] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	panic("Unimplemented")
}

func NewUserHandler(e *echo.Echo, userService service.UserServiceInterface) UserHandlerInterface {
	return &userHandler{userService: userService}
}
