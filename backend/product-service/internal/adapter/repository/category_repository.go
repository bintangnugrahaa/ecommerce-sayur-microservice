package repository

import (
	"context"
	"product-service/internal/core/domain/entity"

	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.CategoryEntity, error)
}

type categoryRepository struct {
	db *gorm.DB
}

// GetAll implements CategoryRepositoryInterface.
func (c *categoryRepository) GetAll(ctx context.Context, queryString entity.QueryStringEntity) ([]entity.CategoryEntity, error) {
	panic("unimplemented")
}

func NewCategoryRepository(db *gorm.DB) CategoryRepositoryInterface {
	return &categoryRepository{db: db}
}
