package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"product-service/internal/core/domain/entity"
	"product-service/internal/core/domain/model"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error)
}

type categoryRepository struct {
	db *gorm.DB
}

// GetAll implements CategoryRepositoryInterface.
func (c *categoryRepository) GetAll(ctx context.Context, query entity.QueryStringEntity) ([]entity.CategoryEntity, int64, int64, error) {
	modelCategories := []model.Category{}
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := c.db.Preload("Products").
		Where("name ILIKE ? OR slug ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	if err := sqlMain.Model(&modelCategories).Count(&countData).Error; err != nil {
		log.Errorf("[CategoryRepository-1] GetAll: %v", err)
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order(order).Limit(int(query.Limit)).Offset(int(offset)).Find(&modelCategories).Error; err != nil {
		log.Errorf("[CategoryRepository-2] GetAll: %v", err)
		return nil, 0, 0, err
	}

	if len(modelCategories) == 0 {
		err := errors.New("404")
		log.Infof("[CategoryRepository-3] GetAll: No category found")
		return nil, 0, 0, err
	}

	entities := []entity.CategoryEntity{}
	for _, val := range modelCategories {
		productEntities := []entity.ProductEntity{}
		for _, prd := range val.Products {
			productEntities = append(productEntities, entity.ProductEntity{
				ID:           prd.ID,
				CategorySlug: val.Slug,
				ParentID:     prd.ParentID,
				Name:         prd.Name,
				Image:        prd.Image,
			})
		}
		status := "Published"
		if val.Status == false {
			status = "Unpublished"
		}

		entities = append(entities, entity.CategoryEntity{
			ID:          val.ID,
			ParentID:    val.ParentID,
			Name:        val.Name,
			Icon:        val.Icon,
			Status:      status,
			Slug:        val.Slug,
			Description: val.Description,
			Products:    productEntities,
		})
	}

	return entities, countData, int64(totalPage), nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepositoryInterface {
	return &categoryRepository{db: db}
}
