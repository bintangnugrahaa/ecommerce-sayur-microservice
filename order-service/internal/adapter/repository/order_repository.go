package repository

import (
	"context"
	"errors"
	"math"
	"order-service/internal/core/domain/entity"
	"order-service/internal/core/domain/model"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error)
	GetByID(ctx context.Context, orderID int64) (*entity.OrderEntity, error)
	CreateOrder(ctx context.Context, req entity.OrderEntity) (int64, error)
	EditOrder(ctx context.Context, req entity.OrderEntity) error
	DeleteOrder(ctx context.Context, orderID int64) error

	GetAllPublished(ctx context.Context) ([]entity.OrderEntity, error)
}

type orderRepository struct {
	db *gorm.DB
}

// CreateOrder implements OrderRepositoryInterface.
func (o *orderRepository) CreateOrder(ctx context.Context, req entity.OrderEntity) (int64, error) {
	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		log.Errorf("[OrderRepository] CreateOrder: %v", err)
		return 0, err
	}

	var orderItems []model.OrderItem
	for _, item := range req.OrderItems {
		orderItem := model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	newOrder := model.Order{
		OrderCode:    req.OrderCode,
		BuyerID:      req.BuyerID,
		OrderDate:    orderDate,
		OrderTime:    req.OrderTime,
		Status:       req.Status,
		TotalAmount:  float64(req.TotalAmount),
		ShippingType: req.ShippingType,
		ShippingFee:  float64(req.ShippingFee),
		Remarks:      req.Remarks,
		OrderItems:   orderItems,
	}

	if err := o.db.Create(&newOrder).Error; err != nil {
		log.Errorf("[OrderRepository] CreateOrder: %v", err)
		return 0, err
	}

	return newOrder.ID, nil
}

// DeleteOrder implements OrderRepositoryInterface.
func (o *orderRepository) DeleteOrder(ctx context.Context, orderID int64) error {
	panic("unimplemented")
}

// EditOrder implements OrderRepositoryInterface.
func (o *orderRepository) EditOrder(ctx context.Context, req entity.OrderEntity) error {
	panic("unimplemented")
}

// GetAll implements OrderRepositoryInterface.
func (o *orderRepository) GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.OrderEntity, int64, int64, error) {
	var modelOrders []model.Order
	var countData int64
	offset := (queryString.Page - 1) * queryString.Limit

	sqlMain := o.db.Preload("OrderItems").
		Where("order_code ILIKE ? OR status ILIKE ?", "%"+queryString.Search+"%", "%"+queryString.Status+"%")
	if err := sqlMain.Model(&modelOrders).Count(&countData).Error; err != nil {
		log.Errorf("[OrderRepository-1] GetAll: %v", err)
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(countData) / float64(queryString.Limit)))
	if err := sqlMain.Order("order_date DESC").Limit(int(queryString.Limit)).Offset(int(offset)).Find(&modelOrders).Error; err != nil {
		log.Errorf("[OrderRepository-2] GetAll: %v", err)
		return nil, 0, 0, err
	}

	if len(modelOrders) == 0 {
		err := errors.New("404")
		log.Infof("[OrderRepository-3] GetAll: No order found")
		return nil, 0, 0, err
	}

	entities := []entity.OrderEntity{}
	for _, val := range modelOrders {
		orderItemsEntities := []entity.OrderItemEntity{}
		for _, item := range val.OrderItems {
			orderItemsEntities = append(orderItemsEntities, entity.OrderItemEntity{
				ID:        item.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			})
		}
		entities = append(entities, entity.OrderEntity{
			ID:          val.ID,
			OrderCode:   val.OrderCode,
			Status:      val.Status,
			OrderDate:   val.OrderDate.Format("2006-01-02 15:04:05"),
			TotalAmount: int64(val.TotalAmount),
			OrderItems:  orderItemsEntities,
			BuyerID:     val.BuyerID,
		})
	}

	return entities, countData, int64(totalPage), nil
}

// GetAllPublished implements OrderRepositoryInterface.
func (o *orderRepository) GetAllPublished(ctx context.Context) ([]entity.OrderEntity, error) {
	panic("unimplemented")
}

// GetByID implements OrderRepositoryInterface.
func (o *orderRepository) GetByID(ctx context.Context, orderID int64) (*entity.OrderEntity, error) {
	var modelOrders model.Order

	if err := o.db.Preload("OrderItems").Where("id =?", orderID).First(&modelOrders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New("404")
			log.Infof("[OrderRepository-1] GetByID: Order not found")
			return nil, err
		}
		log.Errorf("[OrderRepository-2] GetByID: %v", err)
		return nil, err
	}

	orderItemsEntities := []entity.OrderItemEntity{}
	for _, item := range modelOrders.OrderItems {
		orderItemsEntities = append(orderItemsEntities, entity.OrderItemEntity{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &entity.OrderEntity{
		ID:           modelOrders.ID,
		OrderCode:    modelOrders.OrderCode,
		Status:       modelOrders.Status,
		BuyerID:      modelOrders.BuyerID,
		OrderDate:    modelOrders.OrderDate.Format("2006-01-02 15:04:05"),
		TotalAmount:  int64(modelOrders.TotalAmount),
		OrderItems:   orderItemsEntities,
		Remarks:      modelOrders.Remarks,
		ShippingType: modelOrders.ShippingType,
		ShippingFee:  int64(modelOrders.ShippingFee),
	}, nil
}

func NewOrderRepository(db *gorm.DB) OrderRepositoryInterface {
	return &orderRepository{db: db}
}
