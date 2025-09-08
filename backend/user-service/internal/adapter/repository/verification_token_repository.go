package repository

import (
	"context"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"

	"gorm.io/gorm"
)

type VerificationTokenRepositoryInterface interface {
	CreateVerificationToken(ctx context.Context, req entity.VerificationTokenEntity) error
}

type verificationTokenRepository struct {
	db *gorm.DB
}

// CreateVerificationToken implements VerificationTokenRepositoryInterface.
func (v *verificationTokenRepository) CreateVerificationToken(ctx context.Context, req entity.VerificationTokenEntity) error {
	modelVerificationToken := model.VerificationToken{
		UserID:    req.UserID,
		Token:     req.Token,
		TokenType: req.TokenType,
	}

	if err := v.db
	panic("unimplemented")
}

func NewVerificationTokenRepository(db *gorm.DB) VerificationTokenRepositoryInterface {
	return &verificationTokenRepository{db: db}
}
