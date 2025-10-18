package handlers

import (
	"product-service/config"
	"product-service/internal/adapter"
	"product-service/internal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ProductHandlerInterface interface {
	GetAllAdmin(c echo.Context) error
}

type productHandler struct {
	service service.ProductServiceInterface
}

// GetAllAdmin implements ProductHandlerInterface.
func (p *productHandler) GetAllAdmin(c echo.Context) error {
	panic("unimplemented")
}

func NewProductHandler(e *echo.Echo, cfg *config.Config, service service.ProductServiceInterface) ProductHandlerInterface {
	product := &productHandler{service: service}

	e.Use(middleware.Recover())
	mid := adapter.NewMiddlewareAdapter(cfg)
	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/products", product.GetAllAdmin)

	return product
}
