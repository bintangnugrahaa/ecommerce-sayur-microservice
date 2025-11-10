package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"order-service/config"
	httpclient "order-service/internal/adapter/http_client"
	"order-service/internal/adapter/repository"
	"order-service/internal/core/domain/entity"
	"strconv"

	"github.com/labstack/gommon/log"
)

type OrderServiceInterface interface {
	GetAll(ctx context.Context, queryString entity.QueryStringEntity, accessToken string) ([]entity.OrderEntity, int64, int64, error)
}

type orderService struct {
	repo       repository.OrderRepositoryInterface
	cfg        *config.Config
	httpClient httpclient.HttpClient
}

// GetAll implements OrderServiceInterface.
func (o *orderService) GetAll(ctx context.Context, queryString entity.QueryStringEntity, accessToken string) ([]entity.OrderEntity, int64, int64, error) {
	results, count, total, err := o.repo.GetAll(ctx, queryString)
	if err != nil {
		log.Errorf("[OrderService-1] GetAll: %v", err)
		return nil, 0, 0, err
	}

	for _, val := range results {
		for _, res := range val.OrderItems {
			baseUrl := fmt.Sprintf("%s/%s", o.cfg.App.ProductServiceUrl, "admin/products/"+strconv.FormatInt(res.ProductID, 10))
			header := map[string]string{
				"Authorization": "Bearer " + accessToken,
				"Accept":        "application/json",
			}
			dataProduct, err := o.httpClient.CallURL("GET", baseUrl, header, nil)
			if err != nil {
				log.Errorf("[OrderService-2] GetAll: %v", err)
				return nil, 0, 0, err
			}

			defer dataProduct.Body.Close()

			body, err := io.ReadAll(dataProduct.Body)
			if err != nil {
				log.Errorf("[OrderService-3] GetAll: %v", err)
				return nil, 0, 0, err
			}

			var productResponse map[string]interface{}
			err = json.Unmarshal(body, &productResponse)
			if err != nil {
				log.Errorf("[OrderService-4] GetAll: %v", err)
				return nil, 0, 0, err
			}

			res.ProductImage = productResponse["product_image"].(string)
		}
	}

	return results, count, total, nil
}

func NewOrderService(repo repository.OrderRepositoryInterface, cfg *config.Config, httpClient httpclient.HttpClient) OrderServiceInterface {
	return &orderService{
		repo:       repo,
		cfg:        cfg,
		httpClient: httpClient,
	}
}
