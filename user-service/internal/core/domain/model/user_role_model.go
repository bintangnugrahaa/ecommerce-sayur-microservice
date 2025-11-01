package model

import "time"

type UserRole struct {
	ID        int `gorm:"primaryKey"`
	RoleID    int `gorm:"index"`
	UserID    int `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type Tabler interface {
	TableName() string
}

func (UserRole) TableName() string {
	return "user_roles"
}
