package adapter

import (
	"net/http"
	"strings"
	"user-service/config"
	"user-service/internal/core/service"

	"github.com/labstack/echo/v4"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
}

type middlewareAdapter struct {
	cfg *config.Config
}

// CheckToken implements MiddlewareAdapterInterface.
func (m *middlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims := jwt.MapClaims{}
			
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
				}
				return []byte(secretKey), nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}
			
			getSession :=
		}
	}
	panic("unimplemented")
}

func NewMiddlewareAdapter(cfg *config.Config, jwtService service.JwtServiceInterface) MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg: cfg,
	}
}
