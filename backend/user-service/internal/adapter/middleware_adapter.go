package adapter

import (
	"net/http"
	"strings"
	"user-service/config"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
}

type middlewareAdapter struct {
	cfg        *config.Config
	jwtService service.JwtServiceInterface
}

// CheckToken implements MiddlewareAdapterInterface.
func (m *middlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			respErr := response.DefaultResponse{}
			redisConn := config.NewRedisClient()
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				log.Errorf("[MiddlewareAdapter-1] CheckToken: %s", "missing or invalid token")
				respErr.Message = "missing or invalid token"
				respErr.Data = nil
				return c.JSON(http.StatusUnauthorized, respErr)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			_, err := m.jwtService.GenerateToken()

			getSession, err := redisConn.Get(c.Request().Context(), tokenString).Result()
			if err != nil || len(getSession) == 0 {
				log.Errorf("[MiddlewareAdapter-2] CheckToken: %s", err.Error())
				respErr.Message = err.Error()
				respErr.Data = nil
				return c.JSON(http.StatusUnauthorized, respErr)
			}

			c.Set("user", getSession)
			return next(c)
		}
	}
}

func NewMiddlewareAdapter(cfg *config.Config, jwtService service.JwtServiceInterface) MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg:        cfg,
		jwtService: jwtService,
	}
}
