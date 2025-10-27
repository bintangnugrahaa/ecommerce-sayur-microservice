package handlers

import (
	"net/http"
	"product-service/config"
	"product-service/internal/adapter"
	"product-service/internal/adapter/handlers/request"
	"product-service/internal/adapter/handlers/response"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/service"
	"product-service/utils/conv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type ProductHandlerInterface interface {
	GetAllAdmin(c echo.Context) error
	GetByIDAdmin(c echo.Context) error
	CreateAdmin(c echo.Context) error
	EditAdmin(c echo.Context) error
	DeleteAdmin(c echo.Context) error

	GetAllHome(c echo.Context) error
	GetAllShop(c echo.Context) error
}

type productHandler struct {
	service service.ProductServiceInterface
}

// GetAllShop implements ProductHandlerInterface.
func (p *productHandler) GetAllShop(c echo.Context) error {
	var (
		resp      = response.DefaultResponseWithPaginations{}
		ctx       = c.Request().Context()
		respLists = []response.ProductHomeListResponse{}
	)

	orderBy := "created_at"
	if c.QueryParam("orderBy") != "" {
		orderBy = c.QueryParam("orderBy")
	}
	orderType := "desc"
	if c.QueryParam("orderType") != "" {
		orderType = c.QueryParam("orderType")
	}
	var page int64 = 1
	if c.QueryParam("page") != "" {
		page, _ = conv.StringToInt64(c.QueryParam("page"))
	}
	var perPage int64 = 5
	if c.QueryParam("limit") != "" {
		perPage, _ = conv.StringToInt64(c.QueryParam("perPage"))
	}

	var startPrice int64 = 0
	var endPrice int64 = 0
	if c.QueryParam("price") != "" {
		price := strings.Split(c.QueryParam("price"), " - ")
		startPrice, _ = conv.StringToInt64(price[0])
		endPrice, _ = conv.StringToInt64(price[1])
	}

	reqEntity := entity.QueryStringProduct{
		OrderBy:    orderBy,
		OrderType:  orderType,
		Page:       int(page),
		Limit:      int(perPage),
		StartPrice: startPrice,
		EndPrice:   endPrice,
	}

	if c.QueryParam("search") != "" {
		reqEntity.Search = c.QueryParam("search")
	}

	results, totalData, totalPage, err := p.service.GetAll(ctx, reqEntity)
	if err != nil {
		log.Errorf("[ProductHandler-1] GetAllShop: %v", err)
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
		respLists = append(respLists, response.ProductHomeListResponse{
			ID:           result.ID,
			ProductName:  result.Name,
			ProductImage: result.Image,
			SalePrice:    int64(result.SalePrice),
			RegulerPrice: int64(result.RegulerPrice),
			CategoryName: result.CategoryName,
		})
	}

	resp.Message = "success"
	resp.Data = respLists
	resp.Pagination = &response.Pagination{
		Page:       page,
		TotalCount: totalData,
		TotalPage:  totalPage,
		PerPage:    perPage,
	}
	return c.JSON(http.StatusOK, resp)
}

// GetAllHome implements ProductHandlerInterface.
func (p *productHandler) GetAllHome(c echo.Context) error {
	var (
		resp      = response.DefaultResponse{}
		ctx       = c.Request().Context()
		respLists = []response.ProductHomeListResponse{}
	)

	orderBy := "created_at"
	orderType := "desc"
	var page int64 = 1
	var perPage int64 = 5

	reqEntity := entity.QueryStringProduct{
		OrderBy:   orderBy,
		OrderType: orderType,
		Page:      int(page),
		Limit:     int(perPage),
	}

	results, _, _, err := p.service.GetAll(ctx, reqEntity)
	if err != nil {
		log.Errorf("[ProductHandler-1] GetAllHome: %v", err)
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
		respLists = append(respLists, response.ProductHomeListResponse{
			ID:           result.ID,
			ProductName:  result.Name,
			ProductImage: result.Image,
			SalePrice:    int64(result.SalePrice),
			RegulerPrice: int64(result.RegulerPrice),
			CategoryName: result.CategoryName,
		})
	}

	resp.Message = "success"
	resp.Data = respLists
	return c.JSON(http.StatusOK, resp)
}

// DeleteAdmin implements ProductHandlerInterface.
func (p *productHandler) DeleteAdmin(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[ProductHandler-1] DeleteAdmin: %s", "data token not found")
		resp.Message = "data token not found"
		resp.Data = nil
		return c.JSON(http.StatusNotFound, resp)
	}

	idStr := c.Param("id")
	if idStr == "" {
		log.Errorf("[ProductHandler-2] DeleteAdmin: %v", "Invalid id")
		resp.Message = "ID is required"
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	id, err := conv.StringToInt64(idStr)
	if err != nil {
		log.Errorf("[ProductHandler-3] DeleteAdmin: %v", err.Error())
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	err = p.service.Delete(ctx, id)
	if err != nil {
		log.Errorf("[ProductHandler-4] DeleteAdmin: %v", err)
		if err.Error() == "404" {
			resp.Message = "Data not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

// EditAdmin implements ProductHandlerInterface.
func (p *productHandler) EditAdmin(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
		req  = request.ProductRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[ProductHandler-1] EditAdmin: %s", "data token not found")
		resp.Message = "data token not found"
		resp.Data = nil
		return c.JSON(http.StatusNotFound, resp)
	}

	idStr := c.Param("id")
	if idStr == "" {
		log.Errorf("[ProductHandler-2] EditAdmin: %v", "Invalid id")
		resp.Message = "ID is required"
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	id, err := conv.StringToInt64(idStr)
	if err != nil {
		log.Errorf("[ProductHandler-3] EditAdmin: %v", err.Error())
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := c.Bind(&req); err != nil {
		log.Errorf("[ProductHandler-4] EditAdmin: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := c.Validate(req); err != nil {
		log.Errorf("[ProductHandler-3] CreateAdmin: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	reqEntity := entity.ProductEntity{
		ID:           id,
		CategorySlug: req.CategorySlug,
		ParentID:     nil,
		Name:         req.ProductName,
		Image:        req.VariantDetail[0].ProductImage,
		Description:  req.ProductDescription,
		RegulerPrice: float64(req.VariantDetail[0].RegulerPrice),
		SalePrice:    float64(req.VariantDetail[0].SalePrice),
		Unit:         req.Unit,
		Weight:       req.VariantDetail[0].Weight,
		Stock:        req.VariantDetail[0].Stock,
		Variant:      req.Variant,
		Status:       req.Status,
	}

	productChilds := []entity.ProductEntity{}
	if len(req.VariantDetail) > 1 {
		for i := 1; i < len(req.VariantDetail); i++ {
			productChilds = append(productChilds, entity.ProductEntity{
				Image:        req.VariantDetail[i].ProductImage,
				RegulerPrice: float64(req.VariantDetail[i].RegulerPrice),
				SalePrice:    float64(req.VariantDetail[i].SalePrice),
				Weight:       req.VariantDetail[i].Weight,
				Stock:        req.VariantDetail[i].Stock,
			})
		}

		reqEntity.Child = productChilds
	}

	err = p.service.Update(ctx, reqEntity)
	if err != nil {
		log.Errorf("[ProductHandler-4] EditAdmin: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusOK, resp)
}

// CreateAdmin implements ProductHandlerInterface.
func (p *productHandler) CreateAdmin(c echo.Context) error {
	var (
		resp = response.DefaultResponse{}
		ctx  = c.Request().Context()
		req  = request.ProductRequest{}
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[ProductHandler-1] CreateAdmin: %s", "data token not found")
		resp.Message = "data token not found"
		resp.Data = nil
		return c.JSON(http.StatusNotFound, resp)
	}

	if err := c.Bind(&req); err != nil {
		log.Errorf("[ProductHandler-2] CreateAdmin: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := c.Validate(req); err != nil {
		log.Errorf("[ProductHandler-3] CreateAdmin: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	reqEntity := entity.ProductEntity{
		CategorySlug: req.CategorySlug,
		ParentID:     nil,
		Name:         req.ProductName,
		Image:        req.VariantDetail[0].ProductImage,
		Description:  req.ProductDescription,
		RegulerPrice: float64(req.VariantDetail[0].RegulerPrice),
		SalePrice:    float64(req.VariantDetail[0].SalePrice),
		Unit:         req.Unit,
		Weight:       req.VariantDetail[0].Weight,
		Stock:        req.VariantDetail[0].Stock,
		Variant:      req.Variant,
		Status:       req.Status,
	}

	productChilds := []entity.ProductEntity{}
	if len(req.VariantDetail) > 1 {
		for i := 1; i < len(req.VariantDetail); i++ {
			productChilds = append(productChilds, entity.ProductEntity{
				Image:        req.VariantDetail[i].ProductImage,
				RegulerPrice: float64(req.VariantDetail[i].RegulerPrice),
				SalePrice:    float64(req.VariantDetail[i].SalePrice),
				Weight:       req.VariantDetail[i].Weight,
				Stock:        req.VariantDetail[i].Stock,
			})
		}

		reqEntity.Child = productChilds
	}

	err := p.service.Create(ctx, reqEntity)
	if err != nil {
		log.Errorf("[ProductHandler-4] CreateAdmin: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp.Message = "success"
	resp.Data = nil
	return c.JSON(http.StatusCreated, resp)
}

// GetByIDAdmin implements ProductHandlerInterface.
func (p *productHandler) GetByIDAdmin(c echo.Context) error {
	var (
		resp        = response.DefaultResponse{}
		ctx         = c.Request().Context()
		respProduct = response.ProductDetailResponse{}
	)

	user := c.Get("user").(string)
	if user == "" {
		log.Errorf("[ProductHandler-1] GetByIDAdmin: %s", "data token not found")
		resp.Message = "data token not found"
		resp.Data = nil
		return c.JSON(http.StatusNotFound, resp)
	}

	idStr := c.Param("id")
	if idStr == "" {
		log.Errorf("[ProductHandler-2] GetByIDAdmin: %v", "Invalid id")
		resp.Message = "ID is required"
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	id, err := conv.StringToInt64(idStr)
	if err != nil {
		log.Errorf("[ProductHandler-3] GetByIDAdmin: %v", err.Error())
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusBadRequest, resp)
	}

	result, err := p.service.GetByID(ctx, id)
	if err != nil {
		log.Errorf("[ProductHandler-4] GetByIDAdmin: %v", err)
		if err.Error() == "404" {
			resp.Message = "Data not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	responseChilds := []response.ProductChildResponse{}
	if len(result.Child) > 0 {
		for _, child := range result.Child {
			responseChilds = append(responseChilds, response.ProductChildResponse{
				ID:           child.ID,
				SalePrice:    int64(child.SalePrice),
				RegulerPrice: int64(child.RegulerPrice),
				Weight:       child.Weight,
				Stock:        child.Stock,
			})
		}
	}

	respProduct = response.ProductDetailResponse{
		ID:            result.ID,
		ProductName:   result.Name,
		ParentID:      conv.Int64PointerToInt64(result.ParentID),
		ProductImage:  result.Image,
		CategoryName:  result.CategoryName,
		ProductStatus: result.Status,
		SalePrice:     int64(result.SalePrice),
		RegulerPrice:  int64(result.RegulerPrice),
		Unit:          result.Unit,
		Weight:        result.Weight,
		Stock:         result.Stock,
		CreatedAt:     result.CreatedAt,
		Child:         responseChilds,
	}

	resp.Message = "success"
	resp.Data = respProduct
	return c.JSON(http.StatusOK, resp)
}

// GetAllAdmin implements ProductHandlerInterface.
func (p *productHandler) GetAllAdmin(c echo.Context) error {
	var (
		resp         = response.DefaultResponseWithPaginations{}
		ctx          = c.Request().Context()
		respProducts = []response.ProductListResponse{}
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
	if perPageStr := c.QueryParam("limit"); perPageStr != "" {
		perPage, _ = conv.StringToInt64(perPageStr)
		if perPage <= 0 {
			perPage = 10
		}
	}

	categorySlug := c.QueryParam("categorySlug")
	startPrice, err := conv.StringToInt64(c.QueryParam("startPrice"))
	if err != nil {
		startPrice = 0
	}

	endPrice, err := conv.StringToInt64(c.QueryParam("endPrice"))
	if err != nil {
		endPrice = 0
	}

	reqEntity := entity.QueryStringProduct{
		Search:       search,
		OrderBy:      orderBy,
		OrderType:    orderType,
		Page:         int(page),
		Limit:        int(perPage),
		CategorySlug: categorySlug,
		StartPrice:   startPrice,
		EndPrice:     endPrice,
	}

	result, totalData, totalPage, err := p.service.GetAll(ctx, reqEntity)
	if err != nil {
		log.Errorf("[ProductHandler-1] GetAll: %v", err)
		if err.Error() == "404" {
			resp.Message = "Data not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	for _, product := range result {
		respProducts = append(respProducts, response.ProductListResponse{
			ID:            product.ID,
			ProductName:   product.Name,
			ParentID:      conv.Int64PointerToInt64(product.ParentID),
			ProductImage:  product.Image,
			CategoryName:  product.CategoryName,
			ProductStatus: product.Status,
			SalePrice:     int64(product.SalePrice),
			CreatedAt:     product.CreatedAt,
		})
	}

	resp.Data = respProducts
	resp.Message = "success"
	resp.Pagination = &response.Pagination{
		Page:       page,
		TotalCount: totalData,
		TotalPage:  totalPage,
		PerPage:    perPage,
	}

	return c.JSON(http.StatusOK, resp)
}

func NewProductHandler(e *echo.Echo, cfg *config.Config, service service.ProductServiceInterface) ProductHandlerInterface {
	product := &productHandler{service: service}

	e.Use(middleware.Recover())

	homeProduct := e.Group("/products")
	homeProduct.GET("/home", product.GetAllHome)
	homeProduct.GET("/shop", product.GetAllShop)

	mid := adapter.NewMiddlewareAdapter(cfg)
	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/products", product.GetAllAdmin)
	adminGroup.POST("/products", product.CreateAdmin)
	adminGroup.GET("/products/:id", product.GetByIDAdmin)
	adminGroup.PUT("/products/:id", product.EditAdmin)
	adminGroup.DELETE("/products/:id", product.DeleteAdmin)

	return product
}
