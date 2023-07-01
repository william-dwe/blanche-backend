package repository

import (
	"gorm.io/gorm"
)

type OrderItemRepository interface {
}
type OrderItemRepositoryConfig struct {
	DB *gorm.DB
}

type orderItemRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderItemRepository(c OrderItemRepositoryConfig) OrderItemRepository {
	return &orderItemRepositoryImpl{
		db: c.DB,
	}
}
