package entity

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID       uint `gorm:"primaryKey"`
	RoleName string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
