package handlers

import (
	"net/http"
	"product-service/config"
	"product-service/internal/adapter/handlers/response"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service"
	"product-service/utils/conv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type CategoryHandlerInterface interface {
	GetAllAdmin(c echo.Context) error
	GetByIDAdmin(c echo.Context) error
	GetBySlugAdmin(c echo.Context) error
}

type categoryHandler struct {
	categoryService service.CategoryServiceInterface
}

// GetBySlugAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) GetBySlugAdmin(c echo.Context) error {
	panic("unimplemented")
}

// GetByIDAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) GetByIDAdmin(c echo.Context) error {
	panic("unimplemented")
}

// GetAllAdmin implements CategoryHandlerInterface.
func (ch *categoryHandler) GetAllAdmin(c echo.Context) error {
	var (
		resp           = response.DefaultResponseWithPaginations{}
		ctx            = c.Request().Context()
		respCategories = []response.CategoryListAdminResponse{}
	)

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

	results, totalData, totalPage, err := ch.categoryService.GetAll(ctx, reqEntity)
	if err != nil {
		log.Errorf("[CategoryHandler-1] Create: %v", err)
		if err.Error() == "404" {
			resp.Message = "Data not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	for _, result := range results {
		respCategories = append(respCategories, response.CategoryListAdminResponse{
			ID:           result.ID,
			Name:         result.Name,
			Icon:         result.Icon,
			Slug:         result.Slug,
			Status:       result.Status,
			TotalProduct: len(result.Products),
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: totalData,
		PerPage:    perPage,
		TotalPage:  totalPage,
	}
	resp.Message = "success"
	resp.Data = respCategories
	resp.Pagination = &pagination
	return c.JSON(http.StatusOK, resp)
}

func NewCategoryHandler(e *echo.Echo, categoryService service.CategoryServiceInterface, cfg *config.Config) CategoryHandlerInterface {
	category := &categoryHandler{categoryService: categoryService}

	adminGroup := e.Group("/admin")
	adminGroup.GET("/categories", category.GetAllAdmin)

	return category
}
