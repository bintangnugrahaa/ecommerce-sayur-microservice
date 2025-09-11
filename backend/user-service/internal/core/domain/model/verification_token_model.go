package model

import "time"

type VerificationToken struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64 `gorm:"user_id,index"`
	Token     string
	TokenType string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	User      User `gorm:"foreignKey:UserID"`
}
