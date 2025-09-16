package handler

import (
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/storage"
	"user-service/internal/core/service"

	"github.com/labstack/echo/v4"
)

type UploadImageInterface interface {
	UploadImage(c echo.Context) error
}

type uploadImage struct {
	storageHandler storage.SupabaseInterface
}

func NewUploadImage(e *echo.Echo, cfg *config.Config, storageHandler storage.SupabaseInterface, jwtService service.JwtServiceInterface) UploadImageInterface {
	res := &uploadImage{
		storageHandler: storageHandler,
	}

	mid := adapter.NewMiddlewareAdapter(cfg, jwtService)
	e.POST("/auth/profile/image-upload", res., mid.CheckToken())
}