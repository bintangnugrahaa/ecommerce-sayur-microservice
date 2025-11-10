package handler

import (
	"order-service/config"
	"order-service/intenal/adapter"
	"order-service/intenal/adapter/handler/response"
	"order-service/intenal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type OrderHandlerInterface interface {
	GetAllAdmin(c echo.Context) error
}

type orderHandler struct {
	orderService service.OrderServiceInterface
}

// GetAllAdmin implements OrderHandlerInterface.
func (o *orderHandler) GetAllAdmin(c echo.Context) error {
		var (
		ctx            = c.Request().Context()
		respOrders = []response.OrderAdminList{}
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[ProductHandler-1] CreateAdmin: %s", "data token not found")
		return c.JSON(http.StatusNotFound, response.ResponseError("data token not found"))
	}

	search := c.QueryParam("search")
	orderBy := "created_at"
	if c.QueryParam("orderBy") != "" {
		orderBy = c.QueryParam("orderBy")
	}

	orderType := "desc"
	if c.QueryParam("orderType") != "" {
		orderType = c.QueryParam("orderType")
	}

	var page int64 = 1
	if pageStr := c.QueryParam("page"); pageStr != "" {
		page, _ = conv.StringToInt64(pageStr)
		if page <= 0 {
			page = 1
		}
	}

	var perPage int64 = 10
	if perPageStr := c.QueryParam("perPage"); perPageStr != "" {
		perPage, _ = conv.StringToInt64(perPageStr)
		if perPage <= 0 {
			perPage = 10
		}
	}

	reqEntity := entity.QueryStringEntity{
		Search:    search,
		OrderBy:   orderBy,
		OrderType: orderType,
		Page:      int(page),
		Limit:     int(perPage),
	}

	panic("unimplemented")
}

func NewOrderHandler(orderService service.OrderServiceInterface, e *echo.Echo, cfg *config.Config) OrderHandlerInterface {
	ordHandler := &orderHandler{orderService: orderService}

	e.Use(middleware.Recover())
	mid := adapter.NewMiddlewareAdapter(cfg)
	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/orders", ordHandler.GetAllAdmin)

	return ordHandler
}
