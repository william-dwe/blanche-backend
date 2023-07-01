package repository

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type UserOrderRepository interface {
	CreateOrder(newOrder entity.UserOrder) (*entity.UserOrder, error)
	FindOrderByOrderCode(orderCode string, userId uint) (*entity.UserOrder, error)
}

type UserOrderRepositoryConfig struct {
	DB *gorm.DB
}

type userOrderRepositoryImpl struct {
	db *gorm.DB
}

func NewUserOrderRepository(c UserOrderRepositoryConfig) UserOrderRepository {
	return &userOrderRepositoryImpl{
		db: c.DB,
	}
}

func (r *userOrderRepositoryImpl) CreateOrder(newOrder entity.UserOrder) (*entity.UserOrder, error) {
	err := r.db.Create(&newOrder).Error
	if err != nil {
		return nil, domain.ErrCreateUserOrder
	}

	return &newOrder, nil
}

func (r *userOrderRepositoryImpl) FindOrderByOrderCode(orderCode string, userId uint) (*entity.UserOrder, error) {
	var order entity.UserOrder

	err := r.db.Unscoped().
		Where("order_code = ? AND user_id = ?", orderCode, userId).
		Preload("OrderItems").
		Preload("OrderItems.Product", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Preload("OrderItems.Product.Merchant").
		Preload("OrderItems.Product.ProductImages").
		Preload("OrderItems.VariantItem").
		Preload("OrderItems.Product.ProductPromotion").
		Preload("OrderItems.Product.ProductPromotion.Promotion", func(db *gorm.DB) *gorm.DB {
			return db.Where("start_at <= ?", time.Now()).Where("end_at >= ?", time.Now()).Where("quota > 0")
		}).
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetOrderNotFound
		}
		return nil, domain.ErrGetOrderSummary
	}

	return &order, nil
}
